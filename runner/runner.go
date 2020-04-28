package runner

import "flowflux/nodecollection"

// Runner ...
type Runner interface {
	Node() nodecollection.Node
	Input() chan<- InputMessage
	Start()
}
