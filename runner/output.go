package runner

import (
	"bufio"
	"flowflux/nodecollection"
	"log"
	"os"
)

// OutputRunner ...
type OutputRunner struct {
	node              nodecollection.Node
	findOutputRunners func(Runner) []Runner
	channel           chan InputMessage
	didCloseOutput    chan<- bool
}

// Node ...
func (o OutputRunner) Node() nodecollection.Node { return o.node }

// Start ...
func (o OutputRunner) Start() {
	writer := bufio.NewWriter(os.Stdout)

	for message := range o.channel {
		if message.EOF {
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
