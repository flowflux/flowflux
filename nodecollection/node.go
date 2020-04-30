package nodecollection

import (
	"fmt"
	"strings"
)

// Class ...
type Class int

// Classes ...
const (
	UnknownClass Class = iota
	InputClass
	ProcessClass
	ForkClass
	MergeClass
	PipeClass
	OutputClass
)

// ClassToString ...
func ClassToString(class Class) string {
	switch class {
	case InputClass:
		return "input"
	case ProcessClass:
		return "process"
	case ForkClass:
		return "fork"
	case MergeClass:
		return "merge"
	case PipeClass:
		return "pipe"
	case OutputClass:
		return "output"
	default:
		return "unknown"
	}
}

// Node ...
type Node struct {
	Class      Class
	ID         string
	ScanMethod ScanMethod
	Process    ProcessCommand
	OutKeys    []string
}

// ProcessCommand ...
type ProcessCommand struct {
	Command   string
	Arguments []string
	Scaling   uint
}

func (p ProcessCommand) String() string {
	return fmt.Sprintf(
		"%s %s",
		p.Command,
		strings.Join(p.Arguments, " "),
	)
}

// ScanMethod ...
type ScanMethod int

// ScanMethods ...
const (
	ScanMessages ScanMethod = iota
	ScanRawBytes
)
