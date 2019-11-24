package config

import (
	"fmt"

	"github.com/jaqmol/approx/project"
)

// Formation ...
type Formation struct {
	Actors   map[string]Actor
	FlowTree *FlowTree
}

// NewFormation ...
func NewFormation(projForm *project.Formation) (*Formation, error) {
	actrs := make(map[string]Actor, len(projForm.Definitions))
	for name, def := range projForm.Definitions {
		switch def.Type() {
		case project.StdinType:
			actrs[name] = Stdin
		case project.CommandType:
			prCmd := def.(*project.Command)
			actrs[name] = &Command{
				Ident: prCmd.Ident(),
				Cmd:   prCmd.Cmd(),
				Env:   joinKeyValues(prCmd.Env()),
			}
		case project.ForkType:
			actrs[name] = &Fork{
				Ident: def.Ident(),
			}
		case project.MergeType:
			actrs[name] = &Merge{
				Ident: def.Ident(),
			}
		case project.StdoutType:
			actrs[name] = Stdout
		}
	}
	tree, err := NewFlowTree(projForm.Flows, actrs)
	if err != nil {
		return nil, err
	}
	tree.Iterate(func(prev []*FlowNode, curr *FlowNode, next []*FlowNode) error {
		if curr.Actor().Type() == ForkType {
			fork := curr.Actor().(*Fork)
			fork.Count = len(next)
		} else if curr.Actor().Type() == MergeType {
			merge := curr.Actor().(*Merge)
			merge.Count = len(prev)
		}
		return nil
	})
	return &Formation{actrs, tree}, nil
}

func joinKeyValues(mapping map[string]string) []string {
	acc := make([]string, len(mapping))
	idx := 0
	for key, value := range mapping {
		joining := fmt.Sprintf("%v=%v", key, value)
		acc[idx] = joining
		idx++
	}
	return acc
}
