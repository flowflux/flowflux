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
	inputEOF := make(chan bool)

	for _, id := range collection.IDs() {
		n, _ := collection.Node(id)
		switch n.Class {
		case InputClass:
			i := InputRunner{
				node:              n,
				findOutputRunners: findOutputRunners,
				inputEOF:          inputEOF,
			}
			index[n.ID] = i
		case OutputClass:
			o := OutputRunner{
				node:              n,
				findOutputRunners: findOutputRunners,
				channel:           make(chan []byte, channelBufferSize),
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
				channel:           make(chan []byte, channelBufferSize),
			}
			index[n.ID] = i
		case ProcessClass:
			c := ProcessRunner{
				node:              n,
				findOutputRunners: findOutputRunners,
				channel:           make(chan []byte, channelBufferSize),
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

	concludeTimer := time.NewTimer(concludeTimeoutDuration)
	concludeTimer.Stop()
	conclude := func() {
		os.Exit(0)
	}

	for {
		select {
		case msg := <-processErrorMsgs:
			printErrLn(string(msg))
		case <-inputEOF:
			concludeTimer.Stop()
			concludeTimer.Reset(concludeTimeoutDuration)
		case <-concludeTimer.C:
			conclude()
		}
	}

	// for message := range processErrorMsgs {
	// 	printErrLn(string(message))
	// }
}

// Runner ...
type Runner interface {
	Node() Node
	Input() chan<- []byte
	Start()
}

// InfrastructureRunner ...
type InfrastructureRunner struct {
	node              Node
	findOutputRunners func(Runner) []Runner
	channel           chan []byte
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
func (i InfrastructureRunner) Input() chan<- []byte { return i.channel }

// ProcessRunner ...
type ProcessRunner struct {
	node              Node
	findOutputRunners func(Runner) []Runner
	channel           chan []byte
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
		printLogLn(fmt.Sprintf(
			"%s %s will dispatch to %v channels",
			p.node.Process.Command,
			p.node.Process.Arguments,
			len(dispatchChannels),
		))
		scanner := NewHeavyDutyScanner(cmdOut, MsgDelimiter)
		// scanner.Decode = DecodeBase64Message NOT NEEDED DecodeMessage never called
		for scanner.Scan() {
			for _, c := range dispatchChannels {
				c <- scanner.DelimitedMessage()
				printLogLn(fmt.Sprintf("DID DISPATCH LEN %v MSG", len(scanner.Message())))
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
		cmdIn.Write(message)
	}
}

// Input ...
func (p ProcessRunner) Input() chan<- []byte { return p.channel }

// InputRunner ...
type InputRunner struct {
	node              Node
	findOutputRunners func(Runner) []Runner
	channel           chan []byte
	inputEOF          chan<- bool
}

// Node ...
func (i InputRunner) Node() Node { return i.node }

// Start ...
func (i InputRunner) Start() {
	reader := bufio.NewReader(os.Stdin)
	dispatchChannels := collectInputChannels(i.findOutputRunners(i))

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatalf("Error reading from STDIN: %s", err)
		}

		message := make([]byte, len(line))
		copy(message, line)

		for _, c := range dispatchChannels {
			c <- message
		}

		if err == io.EOF {
			i.inputEOF <- true
			break
		}
	}
}

// Input ...
func (i InputRunner) Input() chan<- []byte { return i.channel }

// OutputRunner ...
type OutputRunner struct {
	node              Node
	findOutputRunners func(Runner) []Runner
	channel           chan []byte
}

// Node ...
func (o OutputRunner) Node() Node { return o.node }

// Start ...
func (o OutputRunner) Start() {
	writer := bufio.NewWriter(os.Stdout)

	for message := range o.channel {
		n, err := writer.Write(message)
		if err != nil {
			log.Fatalf("Error writing to STDOUT: %s", err)
		}
		printLogLn(fmt.Sprintf("DID WRITE %v/%v TO STDOUT", n, len(message)))
	}
}

// Input ...
func (o OutputRunner) Input() chan<- []byte { return o.channel }

// collectInputChannels
func collectInputChannels(runners []Runner) []chan<- []byte {
	channels := make([]chan<- []byte, len(runners))
	for i, r := range runners {
		channels[i] = r.Input()
	}
	return channels
}
