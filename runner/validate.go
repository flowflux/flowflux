package runner

import (
	"flowflux/nodecollection"
	"fmt"
)

// ValidateCollection ...
func ValidateCollection(collection nodecollection.Collection) error {
	makeError := func(node nodecollection.Node) error {
		return fmt.Errorf(
			"Error interpreting node of type \"%s\" [%s]",
			nodecollection.ClassToString(node.Class),
			node.ID,
		)
	}
	for _, nodeID := range collection.IDs() {
		node, nodeExists := collection.Node(nodeID)
		if !nodeExists {
			return makeError(node)
		}
		for _, nextNodeID := range node.OutKeys {
			nextNode, nextNodeExists := collection.Node(nextNodeID)
			if !nextNodeExists {
				return makeError(nextNode)
			}
		}
	}
	return nil
}
