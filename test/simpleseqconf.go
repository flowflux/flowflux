package test

import "github.com/jaqmol/approx/config"

// SimpleSequenceConfig ...
type SimpleSequenceConfig struct {
	Fork             config.Fork
	FirstNameExtract config.Command
	LastNameExtract  config.Command
	Merge            config.Merge
}

// MakeSimpleSequenceConfig ...
func MakeSimpleSequenceConfig() *SimpleSequenceConfig {
	mergeConf := config.Merge{
		Ident: "merge",
		Count: 2,
	}
	firstNameExtractConf := config.Command{
		Ident: "extract-first-name",
		Cmd:   "node ../test/node-procs/test-extract-prop.js",
		Env:   []string{"PROP_NAME=first_name"},
	}
	lastNameExtractConf := config.Command{
		Ident: "extract-last-name",
		Cmd:   "node ../test/node-procs/test-extract-prop.js",
		Env:   []string{"PROP_NAME=last_name"},
	}
	forkConf := config.Fork{
		Ident: "fork",
		Count: 2,
	}
	return &SimpleSequenceConfig{
		Fork:             forkConf,
		FirstNameExtract: firstNameExtractConf,
		LastNameExtract:  lastNameExtractConf,
		Merge:            mergeConf,
	}
}
