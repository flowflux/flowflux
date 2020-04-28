package nodecollection

// Collection ...
type Collection struct {
	index     map[string]Node
	indexKeys []string
}

// NewCollection ...
func NewCollection(filePath string) Collection {
	index := parseHubFile(filePath)
	return Collection{
		index: index,
	}
}

// IDs ...
func (c Collection) IDs() []string {
	if c.indexKeys == nil {
		c.indexKeys = make([]string, len(c.index))
		i := 0
		for k := range c.index {
			c.indexKeys[i] = k
			i++
		}
	}
	return c.indexKeys
}

// Node ...
func (c Collection) Node(id string) (Node, bool) {
	n, ok := c.index[id]
	return n, ok
}

// Outputs ...
func (c Collection) Outputs(n Node) []Node {
	nodes := make([]Node, len(n.OutKeys))
	for i, key := range n.OutKeys {
		nodes[i] = c.index[key]
	}
	return nodes
}
