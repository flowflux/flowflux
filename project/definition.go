package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// DefinitionType ...
type DefinitionType int

// DefinitionType ...
const (
	StdinType DefinitionType = iota
	CommandType
	ForkType
	MergeType
	StdoutType
)

// Definition ...
type Definition interface {
	Type() DefinitionType
	Ident() string
}

// LoadDefinition ...
func LoadDefinition(projectDirectory string, flow []Flow) (map[string]Definition, error) {
	definitionFilepath := filepath.Join(projectDirectory, "definition.yaml")
	_, err := os.Stat(definitionFilepath)
	if !os.IsNotExist(err) {
		return loadDefinitionFromPath(definitionFilepath, flow)
	}
	return nil, err
}

func loadDefinitionFromPath(definitionFilepath string, flow []Flow) (map[string]Definition, error) {
	data, err := ioutil.ReadFile(definitionFilepath)
	if err != nil {
		return nil, err
	}
	var parsed map[string]map[string]interface{}
	err = yaml.Unmarshal(data, &parsed)
	if err != nil {
		return nil, err
	}
	return interpreteDefinition(parsed, flow)
}

func interpreteDefinition(dataMap map[string]map[string]interface{}, flow []Flow) (map[string]Definition, error) {
	defs := make(map[string]Definition)
	for name, data := range dataMap {
		switch data["type"] {
		case "command":
			defs[name] = NewCommand(name, data)
		case "fork":
			defs[name] = NewFork(name, data)
		case "merge":
			defs[name] = NewMerge(name, data)
		}
	}
	for _, line := range flow {
		for _, name := range line {
			switch name {
			case "<stdin>":
				_, ok := defs[name]
				if ok {
					return nil, fmt.Errorf("Flow connects <stdin> more than once")
				}
				defs[name] = &Stdin
			case "<stdout>":
				_, ok := defs[name]
				if ok {
					return nil, fmt.Errorf("Flow connects <stdout> more than once")
				}
				defs[name] = &Stdout
			}
		}
	}
	return defs, nil
}
