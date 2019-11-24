package project_test

import (
	"path/filepath"
	"testing"

	"github.com/jaqmol/approx/project"
)

// TestComplexProjectFormation ...
func TestComplexProjectFormation(t *testing.T) {
	// t.SkipNow()                                            // TODO CONTINUE
	projDir, err := filepath.Abs("../test/beta-test-proj") // /formation.yaml
	if err != nil {
		t.Fatal(err)
	}
	form, err := project.LoadFormation(projDir)
	if err != nil {
		t.Fatal(err)
	}

	checkProjectDefinitions(t, form.Definitions, true)
	checkProjectFlows(t, form.Flows, true)
}

// TestSimpleProjectFormation ...
func TestSimpleProjectFormation(t *testing.T) {
	// t.SkipNow()                                     // TODO CONTINUE
	projDir, err := filepath.Abs("../test/alpha-test-proj") // /definition.yaml /flow.yaml
	if err != nil {
		t.Fatal(err)
	}
	form, err := project.LoadFormation(projDir)
	if err != nil {
		t.Fatal(err)
	}

	checkProjectDefinitions(t, form.Definitions, false)
	checkProjectFlows(t, form.Flows, false)
}
