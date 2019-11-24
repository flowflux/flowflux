package project

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Flow ...
type Flow []string

// LoadFlow ...
func LoadFlow(projectDirectory string) ([]Flow, error) {
	flowFilepath := filepath.Join(projectDirectory, "flow.yaml")
	_, err := os.Stat(flowFilepath)
	if !os.IsNotExist(err) {
		return loadFlowFromPath(flowFilepath)
	}
	return nil, err
}

func loadFlowFromPath(flowFilepath string) ([]Flow, error) {
	data, err := ioutil.ReadFile(flowFilepath)
	if err != nil {
		return nil, err
	}
	var flowData [][]string
	err = yaml.Unmarshal(data, &flowData)
	if err != nil {
		return nil, err
	}
	return interpreteFlow(flowData)
}

func interpreteFlow(data [][]string) ([]Flow, error) {
	fl := make([]Flow, len(data))
	for i, f := range data {
		fl[i] = f
	}
	return fl, nil
}
