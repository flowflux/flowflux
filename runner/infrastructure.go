package runner

import "flowflux/nodecollection"

// InfrastructureRunner ...
type InfrastructureRunner struct {
	node                    nodecollection.Node
	collectDispatchChannels func(runner Runner) []chan<- InputMessage
	channel                 chan InputMessage
}

// NewInfrastructureRunner ...
func NewInfrastructureRunner(
	node nodecollection.Node,
	collectDispatchChannels func(runner Runner) []chan<- InputMessage,
) InfrastructureRunner {
	return InfrastructureRunner{
		node:                    node,
		collectDispatchChannels: collectDispatchChannels,
		channel:                 make(chan InputMessage, channelBufferSize),
	}
}

// Node ...
func (i InfrastructureRunner) Node() nodecollection.Node { return i.node }

// Start ...
func (i InfrastructureRunner) Start() {
	dispatchChannels := i.collectDispatchChannels(i)

	for message := range i.channel {
		for _, c := range dispatchChannels {
			c <- message
		}
	}
}

// Input ...
func (i InfrastructureRunner) Input() chan<- InputMessage { return i.channel }
