package event

import (
	"encoding/json"
	"fmt"
)

// LogMsg ...
type LogMsg []string

// UnmarshalLogMsg ...
func UnmarshalLogMsg(b []byte) (LogMsg, error) {
	var lm []string
	err := json.Unmarshal(b, &lm)
	if err != nil {
		return nil, err
	}
	if len(lm) != 2 {
		return nil, fmt.Errorf("Incompatible log-message: %v", lm)
	}
	return lm, nil
}

// Marshal ...
func (lm LogMsg) Marshal() ([]byte, error) {
	return json.Marshal(lm)
}

// IsError ...
func (lm LogMsg) IsError() bool {
	return lm[0] == "error"
}

// IsMsg ...
func (lm LogMsg) IsMsg() bool {
	t := lm[0]
	return t == "info" || t == "debug" || t == "warn"
}

// Error ...
func (lm LogMsg) Error() (*Error, error) {
	if lm.IsError() {
		pb := []byte(lm.Payload())
		return UnmarshalError(pb)
	}
	return nil, nil
}

// Type ...
func (lm LogMsg) Type() string {
	return lm[0]
}

// Payload ...
func (lm LogMsg) Payload() string {
	return lm[1]
}

// PayloadOrError ...
func (lm LogMsg) PayloadOrError() (*string, *Error, error) {
	if lm.IsError() {
		parsed, err := lm.Error()
		if err != nil {
			return nil, nil, err
		}
		return nil, parsed, nil
	} else if lm.IsMsg() {
		p := lm.Payload()
		return &p, nil, nil
	}
	return nil, nil, fmt.Errorf("Unknown log-message: %v", lm)
}
