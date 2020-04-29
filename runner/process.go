package runner

import (
	"bufio"
	"flowflux/flowscan"
	"flowflux/load"
	"flowflux/nodecollection"
	"io"
	"log"
	"os/exec"
	"time"
)

// ProcessRunner ...
type ProcessRunner struct {
	node              nodecollection.Node
	findOutputRunners func(Runner) []Runner
	channel           chan InputMessage
	processErrorMsgs  chan<- []byte
	load              uint64
}

// Node ...
func (p ProcessRunner) Node() nodecollection.Node { return p.node }

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

	// cmd.Process.Pid

	go func() {
		scanner := bufio.NewScanner(cmdErr)
		for scanner.Scan() {
			message := scanner.Bytes()
			p.processErrorMsgs <- message
		}
	}()

	go func() {
		dispatchChannels := collectInputChannels(p.findOutputRunners(p))
		var scanner flowscan.Scanner
		var scannedMessage func() []byte

		if p.node.ScanMethod == nodecollection.ScanMessages {
			dutyScanner := flowscan.NewHeavyDuty(cmdOut, flowscan.MsgDelimiter)
			scannedMessage = dutyScanner.DelimitedMessage
			scanner = dutyScanner

		} else if p.node.ScanMethod == nodecollection.ScanRawBytes {
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

	// FUTURE FEATURE CODE TO OBSERVE CPU-LOAD OF PROCESS
	// AUTO-SCALING

	var loadSample <-chan uint64
	if p.node.Process.AutoScale {
		loadSample = load.StartSampling(cmd.Process.Pid, 1*time.Second)
	}

	for {
		select {
		case message := <-p.channel:
			if message.EOF {
				cmdIn.Close()
			} else {
				_, err := cmdIn.Write(message.payload)
				if err != nil {
					log.Fatalf(
						"Error writing to stdin of \"%s\": %s",
						p.node.Process.String(),
						err,
					)
				}
			}
		case sample := <-loadSample:
			p.load = sample
		}
	}

	// for message := range p.channel {
	// 	if message.EOF {
	// 		cmdIn.Close()
	// 	} else {
	// 		_, err := cmdIn.Write(message.payload)
	// 		if err != nil {
	// 			log.Fatalf(
	// 				"Error writing to stdin of %s %s: %s",
	// 				p.node.Process.Command,
	// 				strings.Join(p.node.Process.Arguments, ", "),
	// 				err,
	// 			)
	// 		}
	// 	}
	// }
}

// Load ...
func (p ProcessRunner) Load() uint64 { return p.load }

// Input ...
func (p ProcessRunner) Input() chan<- InputMessage { return p.channel }
