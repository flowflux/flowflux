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
	// TODO
	// shouldScale := hasScalingNode(collection)
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
	// TODO
	// var statistics load.Statistics

	for _, id := range collection.IDs() {
		n, _ := collection.Node(id)
		switch n.Class {
		case nodecollection.InputClass:
			i := InputRunner{
				node:              n,
				findOutputRunners: findOutputRunners,
			}
			index[n.ID] = i
		case nodecollection.OutputClass:
			o := OutputRunner{
				node:              n,
				findOutputRunners: findOutputRunners,
				channel:           make(chan InputMessage, channelBufferSize),
				didCloseOutput:    didCloseOutput,
			}
			index[n.ID] = o
		case nodecollection.ForkClass:
			fallthrough
		case nodecollection.MergeClass:
			fallthrough
		case nodecollection.PipeClass:
			i := InfrastructureRunner{
				node:              n,
				findOutputRunners: findOutputRunners,
				channel:           make(chan InputMessage, channelBufferSize),
			}
			index[n.ID] = i
		case nodecollection.ProcessClass:
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

	// TODO
	// if shouldScale {
	// 	statistics = load.NewStatistics()
	// 	procs := make([]load.SampleProcess, 0)
	// 	for _, r := range index {
	// 		switch r.Node().Class {
	// 		case nodecollection.ProcessClass:
	// 			pr, ok := r.(ProcessRunner)
	// 			if ok {
	// 				procs = append(procs, pr)
	// 			}
	// 		}
	// 	}
	// 	go statistics.Start(procs)
	// }

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

// TODO
// func hasScalingNode(collection nodecollection.Collection) bool {
// 	for _, id := range collection.IDs() {
// 		n, _ := collection.Node(id)
// 		if n.Process.AutoScale {
// 			return true
// 		}
// 	}
// 	return false
// }
