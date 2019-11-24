package formation

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jaqmol/approx/logging"

	"github.com/jaqmol/approx/actor"
	"github.com/jaqmol/approx/config"
	"github.com/jaqmol/approx/project"
)

const actorInboxSize = 10

// Formation ...
type Formation struct {
	projPath string
	conf     *config.Formation
	finished chan bool
	actables map[string]actor.Actable
	stdin    io.Reader
	stdout   io.Writer
	logger   logging.Logger
}

// NewFormation ...
func NewFormation(
	stdin io.Reader,
	stdout io.Writer,
	logger logging.Logger,
) (*Formation, error) {
	projPath, err := figureProjectPath()
	if err != nil {
		return nil, err
	}

	projForm, err := project.LoadFormation(projPath)
	if err != nil {
		return nil, err
	}

	confForm, err := config.NewFormation(projForm)
	if err != nil {
		return nil, err
	}

	f := Formation{
		projPath: projPath,
		conf:     confForm,
		finished: make(chan bool),
		actables: make(map[string]actor.Actable),
		stdin:    stdin,
		stdout:   stdout,
		logger:   logger,
	}

	err = f.createActables()
	if err != nil {
		return nil, err
	}
	err = f.connectActables()
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func figureProjectPath() (string, error) {
	var projPath string
	var err error
	if len(os.Args) == 2 {
		projPath = os.Args[1]
	} else {
		projPath, err = os.Getwd()
	}
	if err != nil {
		return "", err
	}

	projPath, err = filepath.Abs(projPath)
	if err != nil {
		return "", err
	}

	return projPath, nil
}

// func loadProjectFormation(projPath string) (*project.Formation, error) {
// 	projForm, err := project.LoadFormation(projPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return projForm, nil
// }

func (f *Formation) createActables() error {
	return f.conf.FlowTree.Iterate(func(
		prev []*config.FlowNode,
		curr *config.FlowNode,
		_ []*config.FlowNode,
	) error {
		var err error
		switch curr.Actor().Type() {
		case config.StdinType:
			fallthrough
		case config.MergeType:
			_, err = f.getCreateActable(curr.Actor())
		default:
			if len(prev) != 1 {
				return fmt.Errorf("Expected precisely 1 input for type of processor \"%v\"", curr.Actor().ID())
			}
			_, err = f.getCreateActable(curr.Actor())
		}
		return err
	})
}

func (f *Formation) getCreateActable(currConfProc config.Actor) (actor.Actable, error) {
	id := currConfProc.ID()
	actbl, ok := f.actables[id]
	var err error
	if !ok {
		switch currConfProc.Type() {
		case config.StdinType:
			actbl = newStdinActor(f.stdin)
		case config.CommandType:
			actbl, err = f.newCommandActor(currConfProc.(*config.Command))
		case config.ForkType:
			actbl = f.newForkActor(currConfProc.(*config.Fork))
		case config.MergeType:
			actbl = f.newMergeActor(currConfProc.(*config.Merge))
		case config.StdoutType:
			actbl = newStdoutActor(f.stdout, f.finished)
		}
		f.actables[id] = actbl
	}
	return actbl, err
}

func (f *Formation) newCommandActor(conf *config.Command) (*actor.Command, error) {
	c, err := actor.NewCommandFromConf(actorInboxSize, conf)
	if err != nil {
		return nil, err
	}
	if f.projPath != "" {
		c.Directory(f.projPath)
	}
	f.logger.Add(c.Logging())
	return c, nil
}

func (f *Formation) newForkActor(conf *config.Fork) *actor.Fork {
	return actor.NewFork(actorInboxSize, conf.Ident, conf.Count)
}

func (f *Formation) newMergeActor(conf *config.Merge) *actor.Merge {
	return actor.NewMerge(actorInboxSize, conf.Ident, conf.Count)
}

func (f *Formation) connectActables() error {
	return f.conf.FlowTree.Iterate(func(
		_ []*config.FlowNode,
		curr *config.FlowNode,
		next []*config.FlowNode,
	) error {
		nextActables, err := f.getActables(getNodeIDs(next))
		if err != nil {
			return err
		}
		currID := curr.Actor().ID()
		currActbl, ok := f.actables[currID]
		if !ok {
			return fmt.Errorf("Actor to connect to not found: %v", currID)
		}

		if len(nextActables) > 0 {
			currActbl.Next(nextActables...)
		}
		return nil
	})
}

func getNodeIDs(nodes []*config.FlowNode) []string {
	acc := make([]string, len(nodes))
	for i, n := range nodes {
		acc[i] = n.Actor().ID()
	}
	return acc
}

func (f *Formation) getActables(ids []string) ([]actor.Actable, error) {
	acc := make([]actor.Actable, len(ids))
	for i, id := range ids {
		actbl, ok := f.actables[id]
		if !ok {
			return nil, fmt.Errorf("Could not find \"%v\"", id)
		}
		acc[i] = actbl
	}
	return acc, nil
}

// Start ...
func (f *Formation) Start() <-chan bool {
	go f.logger.Start()
	for _, actbl := range f.actables {
		actbl.Start()
	}
	return f.finished
}
