package runner

import (
	"flowflux/nodecollection"
	"flowflux/printer"
	"time"
)

const channelBufferSize = 50
const concludeTimeoutDuration = 10 * time.Second

// RunCollection ...
func RunCollection(collection nodecollection.Collection) {
	index := make(map[string]Runner)
	collectDispatchChannels := func(runner Runner) []chan<- InputMessage {
		nodes := collection.Outputs(runner.Node())
		channels := make([]chan<- InputMessage, len(nodes))
		for i, node := range nodes {
			dispRunner := index[node.ID]
			channels[i] = dispRunner.Input()
		}
		return channels
	}

	processErrorMsgs := make(chan []byte)
	didCloseOutput := make(chan bool)

	for _, id := range collection.IDs() {
		n, _ := collection.Node(id)
		switch n.Class {
		case nodecollection.InputClass:
			i := NewInputRunner(
				n,
				collectDispatchChannels,
			)
			index[n.ID] = i
		case nodecollection.OutputClass:
			o := NewOutputRunner(
				n,
				didCloseOutput,
			)
			index[n.ID] = o
		case nodecollection.ForkClass:
			fallthrough
		case nodecollection.MergeClass:
			fallthrough
		case nodecollection.PipeClass:
			i := NewInfrastructureRunner(
				n,
				collectDispatchChannels,
			)
			index[n.ID] = i
		case nodecollection.ProcessClass:
			p := NewProcessRunner(
				n,
				collectDispatchChannels,
				processErrorMsgs,
			)
			index[n.ID] = p
		}
	}

	for _, r := range index {
		switch r.Node().Class {
		case nodecollection.ForkClass:
			fallthrough
		case nodecollection.MergeClass:
			fallthrough
		case nodecollection.PipeClass:
			go r.Start()
		}
	}

	for _, r := range index {
		switch r.Node().Class {
		case nodecollection.ProcessClass:
			go r.Start()
		}
	}

	for _, r := range index {
		switch r.Node().Class {
		case nodecollection.InputClass:
			fallthrough
		case nodecollection.OutputClass:
			go r.Start()
		}
	}

	for {
		breakLoop := false
		select {
		case msg := <-processErrorMsgs:
			printer.ErrLn(string(msg))
		case breakLoop = <-didCloseOutput:
		}
		if breakLoop {
			break
		}
	}
}
