package event

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Error ...
type Error struct {
	ColumnNumber int    `json:"columnNumber,omitempty"`
	FileName     string `json:"fileName,omitempty"`
	LineNumber   int    `json:"lineNumber,omitempty"`
	Message      string `json:"message,omitempty"`
	Name         string `json:"name,omitempty"`
	Stack        string `json:"stack,omitempty"`
}

// UnmarshalError ...
func UnmarshalError(b []byte) (*Error, error) {
	var e Error
	err := json.Unmarshal(b, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// LogMsg ...
func (e *Error) LogMsg() (LogMsg, error) {
	data, err := e.Marshal()
	if err != nil {
		return nil, err
	}
	msg := LogMsg{"error", string(data)}
	return msg, nil
}

// Marshal ...
func (e *Error) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Error) Error() string {
	buff := strings.Builder{}

	if e.Name != "" {
		buff.WriteString(e.Name)
	} else {
		buff.WriteString("ERROR")
	}

	if e.ColumnNumber != 0 && e.LineNumber != 0 {
		str := fmt.Sprintf(" @ Ln %v, Col %v", e.ColumnNumber, e.LineNumber)
		buff.WriteString(str)
	}

	if e.FileName != "" {
		str := fmt.Sprintf(" in %v", e.FileName)
		buff.WriteString(str)
	}

	if e.Message != "" {
		str := fmt.Sprintf("\n  %v", e.Message)
		buff.WriteString(str)
	}

	if e.Stack != "" {
		str := fmt.Sprintf("\n  %v", e.Stack)
		buff.WriteString(str)
	}

	buff.WriteString("\n")

	return buff.String()
}
