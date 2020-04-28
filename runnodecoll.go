package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const channelBufferSize = 50
const concludeTimeoutDuration = 10 * time.Second

func validateNodeCollection(collection NodeCollection) error {
	makeError := func(node Node) error {
		return fmt.Errorf(
			"Error interpreting node of type \"%s\" [%s]",
			NodeClassToString(node.Class),
			node.ID,
		)
	}
	for _, nodeID := range collection.IDs() {
		node, nodeExists := collection.Node(nodeID)
		// log.Printf("Checking node id: %s, does exist: %v\n", nodeID, nodeExists)
		if !nodeExists {
			return makeError(node)
		}
		for _, nextNodeID := range node.OutKeys {
			nextNode, nextNodeExists := collection.Node(nextNodeID)
			// log.Printf("  -> connected node id: %s, does exist: %v\n", nextNodeID, nextNodeExists)
			if !nextNodeExists {
				return makeError(nextNode)
			}
		}
	}
	return nil
}

func runNodeCollection(collection NodeCollection) {
	index := make(map[string]Runner)
	findOutputRunners := func(r Runner) []Runner {
		nodes := collection.Outputs(r.Node())
		runners := make([]Runner, len(nodes))
		for i, node := range nodes {
			runners[i] = index[node.ID]
		}
		return runners
	}

	processErrorMsgs := make(chan []byte)
	didCloseOutput := make(chan bool)

	for _, id := range collection.IDs() {
		n, _ := collection.Node(id)
		switch n.Class {
		case InputClass:
			i := InputRunner{
				node:              n,
				findOutputRunners: findOutputRunners,
			}
			index[n.ID] = i
		case OutputClass:
			o := OutputRunner{
				node:              n,
				findOutputRunners: findOutputRunners,
				channel:           make(chan InputMessage, channelBufferSize),
				didCloseOutput:    didCloseOutput,
			}
			index[n.ID] = o
		case ForkClass:
			fallthrough
		case MergeClass:
			fallthrough
		case PipeClass:
			i := InfrastructureRunner{
				node:              n,
				findOutputRunners: findOutputRunners,
				channel:           make(chan InputMessage, channelBufferSize),
			}
			index[n.ID] = i
		case ProcessClass:
			c := ProcessRunner{
				node:              n,
				findOutputRunners: findOutputRunners,
				channel:           make(chan InputMessage, channelBufferSize),
				processErrorMsgs:  processErrorMsgs,
			}
			index[n.ID] = c
		}
	}

	for _, r := range index {
		switch r.Node().Class {
		case ForkClass:
			fallthrough
		case MergeClass:
			fallthrough
		case PipeClass:
			go r.Start()
		}
	}

	for _, r := range index {
		switch r.Node().Class {
		case ProcessClass:
			go r.Start()
		}
	}

	for _, r := range index {
		switch r.Node().Class {
		case InputClass:
			fallthrough
		case OutputClass:
			go r.Start()
		}
	}

	// concludeTimer := time.NewTimer(concludeTimeoutDuration)
	// concludeTimer.Stop()
	// conclude := func() {
	// 	os.Exit(0)
	// }

	for {
		breakLoop := false
		select {
		case msg := <-processErrorMsgs:
			printErrLn(string(msg))
		case breakLoop = <-didCloseOutput:
		}
		if breakLoop {
			break
		}
	}

	// for message := range processErrorMsgs {
	// 	printErrLn(string(message))
	// }
}

// Runner ...
type Runner interface {
	Node() Node
	Input() chan<- InputMessage
	Start()
}

// InputMessage ...
type InputMessage struct {
	payload []byte
	EOF     bool
}

// InfrastructureRunner ...
type InfrastructureRunner struct {
	node              Node
	findOutputRunners func(Runner) []Runner
	channel           chan InputMessage
}

// Node ...
func (i InfrastructureRunner) Node() Node { return i.node }

// Start ...
func (i InfrastructureRunner) Start() {
	dispatchChannels := collectInputChannels(i.findOutputRunners(i))

	for message := range i.channel {
		for _, c := range dispatchChannels {
			c <- message
		}
	}
}

// Input ...
func (i InfrastructureRunner) Input() chan<- InputMessage { return i.channel }

// ProcessRunner ...
type ProcessRunner struct {
	node              Node
	findOutputRunners func(Runner) []Runner
	channel           chan InputMessage
	processErrorMsgs  chan<- []byte
}

// Node ...
func (p ProcessRunner) Node() Node { return p.node }

// Start ...
func (p ProcessRunner) Start() {
	cmd := exec.Command(p.node.Process.Command, p.node.Process.Arguments...)
	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmdIn, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		scanner := bufio.NewScanner(cmdErr)
		for scanner.Scan() {
			message := scanner.Bytes()
			p.processErrorMsgs <- message
		}
	}()

	go func() {
		dispatchChannels := collectInputChannels(p.findOutputRunners(p))
		var scanner FlowScanner
		var scannedMessage func() []byte

		if p.node.ScanMethod == ScanMessages {
			dutyScanner := NewHeavyDutyScanner(cmdOut, MsgDelimiter)
			scannedMessage = dutyScanner.DelimitedMessage
			scanner = dutyScanner

		} else if p.node.ScanMethod == ScanRawBytes {
			bytesScanner := NewRawBytesScanner(cmdOut)
			scannedMessage = bytesScanner.Message
			scanner = bytesScanner
		}

		for scanner.Scan() {
			for _, c := range dispatchChannels {
				c <- InputMessage{
					payload: scannedMessage(),
				}
			}
		}

		if scanner.Err() != nil {
			if scanner.Err() == io.EOF {
				for _, c := range dispatchChannels {
					c <- InputMessage{
						EOF: true,
					}
				}
			} else {
				log.Fatal(scanner.Err())
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	printLogLn(fmt.Sprintf(
		"Did start: %s %s",
		p.node.Process.Command,
		strings.Join(p.node.Process.Arguments, ", "),
	))

	for message := range p.channel {
		if message.EOF {
			printLogLn(fmt.Sprintf(
				"COMMAND PROCESS %s %s DISPATCHING EOF",
				p.node.Process.Command,
				strings.Join(p.node.Process.Arguments, ", "),
			))
			cmdIn.Close()
		} else {
			_, err := cmdIn.Write(message.payload)
			if err != nil {
				log.Fatalf(
					"Error writing to stdin of %s %s: %s",
					p.node.Process.Command,
					strings.Join(p.node.Process.Arguments, ", "),
					err,
				)
			}
		}
	}
}

// Input ...
func (p ProcessRunner) Input() chan<- InputMessage { return p.channel }

// InputRunner ...
type InputRunner struct {
	node              Node
	findOutputRunners func(Runner) []Runner
	channel           chan InputMessage
}

// Node ...
func (i InputRunner) Node() Node { return i.node }

// Start ...
func (i InputRunner) Start() {
	dispatchChannels := collectInputChannels(i.findOutputRunners(i))
	var scanner FlowScanner
	var scannedMessage func() []byte

	if i.node.ScanMethod == ScanMessages {
		dutyScanner := NewHeavyDutyScanner(os.Stdin, MsgDelimiter)
		scannedMessage = dutyScanner.DelimitedMessage
		scanner = dutyScanner

	} else if i.node.ScanMethod == ScanRawBytes {
		bytesScanner := NewRawBytesScanner(os.Stdin)
		scannedMessage = bytesScanner.Message
		scanner = bytesScanner

		// reader := bufio.NewReader(os.Stdin)
		// bufferSize := 1000
		// buffer := make([]byte, bufferSize)
		// dispatchChannels := collectInputChannels(i.findOutputRunners(i))
		//
		// for {
		// 	n, err := reader.Read(buffer)
		// 	if err != nil && err != io.EOF {
		// 		log.Fatalf("Error reading from STDIN: %s", err)
		// 	}
		//
		// 	message := make([]byte, n)
		// 	copy(message, buffer)
		//
		// 	for _, c := range dispatchChannels {
		// 		c <- message
		// 	}
		//
		// 	if err == io.EOF {
		// 		i.inputEOF <- true
		// 		break
		// 	}
		// }
	}

	for scanner.Scan() {
		for _, c := range dispatchChannels {
			c <- InputMessage{
				payload: scannedMessage(),
			}
		}
	}
	if scanner.Err() != nil {
		if scanner.Err() == io.EOF {
			for _, c := range dispatchChannels {
				printLogLn("INPUT DISPATCHING EOF")
				c <- InputMessage{
					EOF: true,
				}
			}
		} else {
			log.Fatal(scanner.Err())
		}
	}
}

// Input ...
func (i InputRunner) Input() chan<- InputMessage { return i.channel }

// OutputRunner ...
type OutputRunner struct {
	node              Node
	findOutputRunners func(Runner) []Runner
	channel           chan InputMessage
	didCloseOutput    chan<- bool
}

// Node ...
func (o OutputRunner) Node() Node { return o.node }

// Start ...
func (o OutputRunner) Start() {
	writer := bufio.NewWriter(os.Stdout)

	for message := range o.channel {
		// printLogLn(fmt.Sprintf("OutputRunner did receive message of length: %v\n", len(message)))
		if message.EOF {
			printLogLn("OUTPUT DISPATCHING EOF")
			os.Stdout.Close()
			o.didCloseOutput <- true
		} else {
			_, err := writer.Write(message.payload)
			if err != nil {
				log.Fatalf("Error writing final flow output to stdout: %s", err)
			}
		}
	}
}

// Input ...
func (o OutputRunner) Input() chan<- InputMessage { return o.channel }

// collectInputChannels
func collectInputChannels(runners []Runner) []chan<- InputMessage {
	channels := make([]chan<- InputMessage, len(runners))
	for i, r := range runners {
		channels[i] = r.Input()
	}
	return channels
}
