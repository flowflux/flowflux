package project_test

import (
	"path/filepath"
	"testing"

	"github.com/jaqmol/approx/project"
)

// TestProjectFlow ...
func TestProjectFlow(t *testing.T) {
	// t.SkipNow()
	projDir, err := filepath.Abs("../test/alpha-test-proj") // /flow.yaml
	if err != nil {
		t.Fatal(err)
	}
	flows, err := project.LoadFlow(projDir)
	if err != nil {
		t.Fatal(err)
	}

	checkProjectFlows(t, flows, false)
}

func checkProjectFlows(t *testing.T, flows []project.Flow, isComplex bool) {
	var expected [][]string

	if isComplex {
		expected = [][]string{
			[]string{"<stdin>", "fork"},
			[]string{"fork", "extract-first-name", "merge"},
			[]string{"fork", "extract-last-name", "merge"},
			[]string{"merge", "<stdout>"},
		}
	} else {
		expected = [][]string{
			[]string{"fork", "extract-first-name", "merge"},
			[]string{"fork", "extract-last-name", "merge"},
		}
	}

	if len(flows) != len(expected) {
		t.Fatalf("Expected \"%v\" flows, but got \"%v\"", len(expected), len(flows))
	}
	if len(flows[0]) != len(expected[0]) {
		t.Fatalf("Expected \"%v\" flow items @ idx 0, but got \"%v\"", len(expected[0]), len(flows[0]))
	}
	if len(flows[1]) != len(expected[1]) {
		t.Fatalf("Expected \"%v\" flow items @ idx 1, but got \"%v\"", len(expected[1]), len(flows[1]))
	}

	for i, f := range flows {
		for j, name := range f {
			expName := expected[i][j]
			if name != expName {
				t.Fatalf("Expected \"%v\", but got \"%v\"", expName, name)
			}
		}
	}
}
