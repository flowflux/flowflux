package config

// FlowNode ...
type FlowNode struct {
	previous []*FlowNode
	next     []*FlowNode
	actor    Actor
}

// NewFlowNode ...
func NewFlowNode(act Actor) *FlowNode {
	return &FlowNode{
		previous: make([]*FlowNode, 0),
		next:     make([]*FlowNode, 0),
		actor:    act,
	}
}

// Previous ...
func (fn *FlowNode) Previous() []*FlowNode {
	return fn.previous
}

// AppendPrevious ...
func (fn *FlowNode) AppendPrevious(previous ...*FlowNode) {
	fn.previous = append(fn.previous, previous...)
}

// Next ...
func (fn *FlowNode) Next() []*FlowNode {
	return fn.next
}

// AppendNext ...
func (fn *FlowNode) AppendNext(next ...*FlowNode) {
	fn.next = append(fn.next, next...)
}

// Actor ...
func (fn *FlowNode) Actor() Actor {
	return fn.actor
}

// Iterate ...
func (fn *FlowNode) Iterate(callback func(prev []*FlowNode, curr *FlowNode, next []*FlowNode) error) (err error) {
	err = callback(fn.previous, fn, fn.next)
	if err != nil {
		return
	}
	for _, next := range fn.next {
		err = next.Iterate(callback)
		if err != nil {
			return
		}
	}
	return
}
