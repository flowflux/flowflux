package config_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/project"
)

// TeTestConfigurationFormation ...
func TestConfigurationFormation(t *testing.T) {
	projDir, err := filepath.Abs("../test/beta-test-proj") // /formation.yaml
	if err != nil {
		t.Fatal(err)
	}
	projForm, err := project.LoadFormation(projDir)
	if err != nil {
		t.Fatal(err)
	}

	confForm, err := config.NewFormation(projForm)
	if err != nil {
		t.Fatal(err)
	}

	checkDeclarationOfActors(t, confForm)
	checkDeclarationOfFlowNode(t, confForm)
}

func checkDeclarationOfActors(t *testing.T, form *config.Formation) {
	var ok bool
	_, ok = form.Actors["<stdin>"]
	if !ok {
		t.Fatal("Expected <stdin> configuration, but none found")
	}
	_, ok = form.Actors["fork"]
	if !ok {
		t.Fatal("Expected fork configuration, but none found")
	}
	_, ok = form.Actors["extract-first-name"]
	if !ok {
		t.Fatal("Expected extract-first-name configuration, but none found")
	}
	_, ok = form.Actors["extract-last-name"]
	if !ok {
		t.Fatal("Expected extract-last-name configuration, but none found")
	}
	_, ok = form.Actors["merge"]
	if !ok {
		t.Fatal("Expected merge configuration, but none found")
	}
	_, ok = form.Actors["<stdout>"]
	if !ok {
		t.Fatal("Expected <stdout> configuration, but none found")
	}
}

func checkDeclarationOfFlowNode(t *testing.T, form *config.Formation) {
	visited := make(map[string]int)

	checkLen := lengthChecker(map[string][]int{
		"<stdin>":            []int{0, 1},
		"fork":               []int{1, 2},
		"extract-first-name": []int{1, 1},
		"extract-last-name":  []int{1, 1},
		"merge":              []int{2, 1},
		"<stdout>":           []int{1, 0},
	})

	err := form.FlowTree.Iterate(func(prev []*config.FlowNode, curr *config.FlowNode, next []*config.FlowNode) error {
		id := curr.Actor().ID()
		visited[id]++
		return checkLen(id, len(prev), len(next))
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(visited) != 6 {
		t.Fatal("Expected to visit 6 nodes, but got:", len(visited))
	}
	errs := checkContainsAllTimes(
		visited,
		map[string]int{
			"<stdin>":            1,
			"fork":               1,
			"extract-first-name": 1,
			"extract-last-name":  1,
			"merge":              1,
			"<stdout>":           1,
		},
	)
	if len(errs) > 0 {
		err := fmt.Errorf("Errors visiting nodes: %v", strings.Join(errorsToStrings(errs), ", "))
		t.Fatal(err)
	}

	if form.FlowTree.Input == nil {
		t.Fatal("Expected node tree to have an input node")
	}
	if form.FlowTree.Output == nil {
		t.Fatal("Expected node tree to have an output node")
	}

	if form.FlowTree.Input.Actor().ID() != "<stdin>" {
		t.Fatalf("Expected node tree input to be <stdin>, but found: %v", form.FlowTree.Input.Actor().ID())
	}
	if form.FlowTree.Output.Actor().ID() != "<stdout>" {
		t.Fatalf("Expected node tree output to be <stdout>, but found: %v", form.FlowTree.Output.Actor().ID())
	}
}
