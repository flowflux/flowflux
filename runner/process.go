package runner

import (
	"bufio"
	"flowflux/flowscan"
	"flowflux/load"
	"flowflux/nodecollection"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

// ProcessRunner ...
type ProcessRunner struct {
	node                    nodecollection.Node
	collectDispatchChannels func(runner Runner) []chan<- InputMessage
	channel                 chan InputMessage
	processErrorMsgs        chan<- []byte
	instances               []processInstance
}

// NewProcessRunner ...
func NewProcessRunner(
	node nodecollection.Node,
	collectDispatchChannels func(runner Runner) []chan<- InputMessage,
	processErrorMsgs chan<- []byte,
) ProcessRunner {
	instances := make([]processInstance, node.Process.Scaling)
	var i uint
	for i = 0; i < node.Process.Scaling; i++ {
		instances[i] = processInstance{}
	}

	return ProcessRunner{
		node:                    node,
		collectDispatchChannels: collectDispatchChannels,
		channel:                 make(chan InputMessage, channelBufferSize),
		processErrorMsgs:        processErrorMsgs,
		instances:               instances,
	}
}

// Node ...
func (p ProcessRunner) Node() nodecollection.Node { return p.node }

// Start ...
func (p ProcessRunner) Start() {
	dispatchChannels := p.collectDispatchChannels(p)
	instanceChannels := make([]chan InputMessage, p.node.Process.Scaling)

	var i uint
	for i = 0; i < p.node.Process.Scaling; i++ {
		instChan := make(chan InputMessage, channelBufferSize)
		instanceChannels[i] = instChan
		go p.instances[i].start(
			p.node,
			dispatchChannels,
			instChan,
			p.processErrorMsgs,
		)
	}

	// Round-Robin
	instChanIdx := 0
	for message := range p.channel {
		nextInstChan := instanceChannels[instChanIdx]
		nextInstChan <- message
		instChanIdx++
		if instChanIdx == len(instanceChannels) {
			instChanIdx = 0
		}
	}
}

// Input ...
func (p ProcessRunner) Input() chan<- InputMessage { return p.channel }

// TakeLoadSamples ...
func (p ProcessRunner) TakeLoadSamples(every time.Duration, channel chan<- load.ProcessSample) {
	for _, inst := range p.instances {
		inst.takeLoadSamples(every, channel)
	}
}

// processInstance ...
type processInstance struct {
	pid int
}

func (i processInstance) start(
	node nodecollection.Node,
	dispatchChannels []chan<- InputMessage,
	channel chan InputMessage,
	processErrorMsgs chan<- []byte,
) {
	cmd := exec.Command(node.Process.Command, node.Process.Arguments...)
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

	// cmd.Process.Pid

	go func() {
		scanner := bufio.NewScanner(cmdErr)
		for scanner.Scan() {
			message := scanner.Bytes()
			processErrorMsgs <- message
		}
	}()

	go func() {
		var scanner flowscan.Scanner
		var scannedMessage func() []byte

		if node.ScanMethod == nodecollection.ScanMessages {
			// dutyScanner := flowscan.NewHeavyDuty(cmdOut, flowscan.MsgDelimiter)
			// scannedMessage = dutyScanner.DelimitedMessage
			// scanner = dutyScanner
			lenScanner := flowscan.NewLengthPrefix(cmdOut)
			scannedMessage = lenScanner.PrefixedMessage
			scanner = lenScanner

		} else if node.ScanMethod == nodecollection.ScanRawBytes {
			bytesScanner := flowscan.NewRawBytes(cmdOut)
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
	i.pid = cmd.Process.Pid

	// for {
	// 	select {
	// 	case message := <-p.channel:
	// 		if message.EOF {
	// 			cmdIn.Close()
	// 		} else {
	// 			_, err := cmdIn.Write(message.payload)
	// 			if err != nil {
	// 				log.Fatalf(
	// 					"Error writing to stdin of \"%s\": %s",
	// 					p.node.Process.String(),
	// 					err,
	// 				)
	// 			}
	// 		}
	// 	case sample := <-loadSample:
	// 		p.load = sample
	// 	}
	// }

	for message := range channel {
		if message.EOF {
			cmdIn.Close()
		} else {
			_, err := cmdIn.Write(message.payload)
			if err != nil {
				log.Fatalf(
					"Error writing to stdin of %s %s: %s",
					node.Process.Command,
					strings.Join(node.Process.Arguments, ", "),
					err,
				)
			}
		}
	}
}

func (i processInstance) takeLoadSamples(every time.Duration, channel chan<- load.ProcessSample) {
	go load.StartSamplingProcess(i.pid, every, channel)
}
