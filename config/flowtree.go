package config

import (
	"fmt"
	"strings"

	"github.com/jaqmol/approx/project"
)

// FlowTree ...
type FlowTree struct {
	Root   *FlowNode
	Input  *FlowNode
	Output *FlowNode
}

// TODO:
// Loopback-applications like a web-server don't connect to stdin and stdout.
// In this case the definition of one command actor must specify as root.

// NewFlowTree ...
func NewFlowTree(flows []project.Flow, actrs map[string]Actor) (*FlowTree, error) {
	nodeForName := make(map[string]*FlowNode)
	for _, line := range flows {
		for j, toName := range line {
			if j > 0 {
				fromName := line[j-1]
				fromNode := getCreateNode(fromName, actrs[fromName], nodeForName)
				toNode := getCreateNode(toName, actrs[toName], nodeForName)
				fromNode.AppendNext(toNode)
				toNode.AppendPrevious(fromNode)
			} else {
				getCreateNode(toName, actrs[toName], nodeForName)
			}
		}
	}

	for _, node := range nodeForName {
		node.previous = makeUniqueSet(node.previous)
		node.next = makeUniqueSet(node.next)
	}

	input, output, err := findNoPredecessorAndNoSuccessorNodes(nodeForName)
	// TODO: ^ intput and output can be nil in case of a loopback app
	if err != nil {
		return nil, err
	}
	return &FlowTree{
		Root:   input, // TODO: Loopback apps need diverging handling
		Input:  input,
		Output: output,
	}, nil
}

// Iterate ...
func (ft *FlowTree) Iterate(callback func(prev []*FlowNode, curr *FlowNode, next []*FlowNode) error) error {
	wasVisitedForID := make(map[string]bool)
	return ft.Root.Iterate(func(prev []*FlowNode, curr *FlowNode, next []*FlowNode) error {
		id := curr.actor.ID()
		_, ok := wasVisitedForID[id]
		if !ok {
			wasVisitedForID[id] = true
			return callback(prev, curr, next)
		}
		return nil
	})
}

func findNoPredecessorAndNoSuccessorNodes(nodeForName map[string]*FlowNode) (
	noPredecessor *FlowNode,
	noSuccessor *FlowNode,
	err error,
) {
	inputNodes := make([]*FlowNode, 0)
	outputNodes := make([]*FlowNode, 0)
	for _, node := range nodeForName {
		if len(node.previous) == 0 {
			inputNodes = append(inputNodes, node)
		}
		if len(node.next) == 0 {
			outputNodes = append(outputNodes, node)
		}
	}
	// Happy path
	if len(inputNodes) == 1 {
		noPredecessor = inputNodes[0]
	}
	if len(outputNodes) == 1 {
		noSuccessor = outputNodes[0]
	}
	if len(inputNodes) < 2 || len(outputNodes) < 2 {
		return
	}
	// Fail path
	inputIds := collectIDs(inputNodes)
	outputIds := collectIDs(outputNodes)

	allErrMsgs := make([]string, 0)
	if len(inputIds) > 0 {
		allErrMsgs = append(allErrMsgs, moreThanOneErrorMsg(inputIds))
	}
	if len(outputIds) > 0 {
		allErrMsgs = append(allErrMsgs, moreThanOneErrorMsg(outputIds))
	}
	err = fmt.Errorf("Error(s) interpreting flow: %v", strings.Join(allErrMsgs, "; "))
	return nil, nil, err
}

func getCreateNode(name string, act Actor, acc map[string]*FlowNode) *FlowNode {
	node, ok := acc[name]
	if !ok {
		node = NewFlowNode(act)
		acc[name] = node
	}
	return node
}

func makeUniqueSet(input []*FlowNode) []*FlowNode {
	output := make([]*FlowNode, 0)
	isContainedForID := make(map[string]bool)
	for _, node := range input {
		id := node.actor.ID()
		_, isOk := isContainedForID[id]
		if !isOk {
			isContainedForID[id] = true
			output = append(output, node)
		}
	}
	return output
}

func collectIDs(nodes []*FlowNode) []string {
	ids := make([]string, len(nodes))
	for i, n := range nodes {
		ids[i] = n.actor.ID()
	}
	return ids
}

func moreThanOneErrorMsg(ids []string) string {
	return fmt.Sprintf("More than 1 input node: %v", strings.Join(ids, ", "))
}
