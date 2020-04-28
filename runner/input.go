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
	node              nodecollection.Node
	findOutputRunners func(Runner) []Runner
	channel           chan InputMessage
}

// Node ...
func (i InputRunner) Node() nodecollection.Node { return i.node }

// Start ...
func (i InputRunner) Start() {
	dispatchChannels := collectInputChannels(i.findOutputRunners(i))
	var scanner flowscan.Scanner
	var scannedMessage func() []byte

	if i.node.ScanMethod == nodecollection.ScanMessages {
		dutyScanner := flowscan.NewHeavyDuty(os.Stdin, flowscan.MsgDelimiter)
		scannedMessage = dutyScanner.DelimitedMessage
		scanner = dutyScanner

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
func (i InputRunner) Input() chan<- InputMessage { return i.channel }
