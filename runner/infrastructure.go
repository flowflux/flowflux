package runner

import "flowflux/nodecollection"

// InfrastructureRunner ...
type InfrastructureRunner struct {
	node              nodecollection.Node
	findOutputRunners func(Runner) []Runner
	channel           chan InputMessage
}

// Node ...
func (i InfrastructureRunner) Node() nodecollection.Node { return i.node }

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
