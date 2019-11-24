package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// "gopkg.in/yaml.v2"

// Formation ...
type Formation struct {
	Definitions map[string]Definition // `yaml:"definition,omitempty"`
	Flows       []Flow                // `yaml:"flow,omitempty"`
}

// LoadFormation ...
func LoadFormation(projectDirectory string) (*Formation, error) {
	formationFilepath := filepath.Join(projectDirectory, "formation.yaml")
	if _, err := os.Stat(formationFilepath); !os.IsNotExist(err) {
		return loadFormationFromPath(formationFilepath)
	}
	errsAcc := make([]error, 0)

	flow, err := LoadFlow(projectDirectory)
	if err != nil {
		errsAcc = append(errsAcc, err)
	}

	def, err := LoadDefinition(projectDirectory, flow)
	if err != nil {
		errsAcc = append(errsAcc, err)
	}

	if len(errsAcc) > 0 {
		return nil, noFormationError(errsAcc)
	}
	f := Formation{
		Definitions: def,
		Flows:       flow,
	}
	return &f, nil
}

func loadFormationFromPath(formationFilepath string) (*Formation, error) {
	data, err := ioutil.ReadFile(formationFilepath)
	if err != nil {
		return nil, err
	}
	var forMap map[string]interface{}
	err = yaml.Unmarshal(data, &forMap)
	if err != nil {
		return nil, err
	}

	flowData := toListListStr(forMap["flow"])
	flows, err := interpreteFlow(flowData)
	if err != nil {
		return nil, err
	}

	defData := toMapStrMapStrIf(forMap["definition"])
	defs, err := interpreteDefinition(defData, flows)
	if err != nil {
		return nil, err
	}

	return &Formation{defs, flows}, nil
}

func noFormationError(causes []error) error {
	buff := strings.Builder{}
	for i, e := range causes {
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(e.Error())
	}
	return fmt.Errorf("No formation.yaml found, expected definition.yaml and flow.yaml, but found error(s): %v", buff.String())
}

func toMapStrMapStrIf(originalData interface{}) map[string]map[string]interface{} {
	dataList := originalData.(map[interface{}]interface{})
	rootAcc := make(map[string]map[string]interface{})
	for ifName, ifData := range dataList {
		mapAcc := make(map[string]interface{})
		for ifKey, ifValue := range ifData.(map[interface{}]interface{}) {
			mapAcc[ifKey.(string)] = ifValue
		}
		rootAcc[ifName.(string)] = mapAcc
	}
	return rootAcc
}

func toListListStr(originalData interface{}) [][]string {
	dataList := originalData.([]interface{})
	acc := make([][]string, len(dataList))
	for i, ifData := range dataList {
		ifList := ifData.([]interface{})
		strList := make([]string, len(ifList))
		for j, ifData := range ifList {
			strList[j] = ifData.(string)
		}
		acc[i] = strList
	}
	return acc
}
