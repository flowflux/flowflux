package runner

import (
	"flowflux/flowscan"
	"flowflux/nodecollection"
	"io"
	"log"
	"os"
)

// InputRunner ...
type InputRunner struct {
	node                    nodecollection.Node
	collectDispatchChannels func(runner Runner) []chan<- InputMessage
}

// NewInputRunner ...
func NewInputRunner(
	node nodecollection.Node,
	collectDispatchChannels func(runner Runner) []chan<- InputMessage,
) InputRunner {
	return InputRunner{
		node:                    node,
		collectDispatchChannels: collectDispatchChannels,
	}
}

// Node ...
func (i InputRunner) Node() nodecollection.Node { return i.node }

// Start ...
func (i InputRunner) Start() {
	dispatchChannels := i.collectDispatchChannels(i)
	var scanner flowscan.Scanner
	var scannedMessage func() []byte

	if i.node.ScanMethod == nodecollection.ScanMessages {
		// dutyScanner := flowscan.NewHeavyDuty(os.Stdin, flowscan.MsgDelimiter)
		// scannedMessage = dutyScanner.DelimitedMessage
		// scanner = dutyScanner
		lenScanner := flowscan.NewLengthPrefix(os.Stdin)
		scannedMessage = lenScanner.PrefixedMessage
		scanner = lenScanner

	} else if i.node.ScanMethod == nodecollection.ScanRawBytes {
		bytesScanner := flowscan.NewRawBytes(os.Stdin)
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
}

// Input ...
func (i InputRunner) Input() chan<- InputMessage {
	log.Fatalln("Flow input cannot be used as output of other process")
	return nil
}
