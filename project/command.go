package project

// Command ...
type Command struct {
	ident string            // `yaml:"ident,omitempty"`
	cmd   string            // `yaml:"cmd,omitempty"`
	env   map[string]string // `yaml:"env,omitempty"`
}

// NewCommand ...
func NewCommand(ident string, originalData interface{}) *Command {
	data := originalData.(map[string]interface{})
	c := Command{
		ident: ident,
		cmd:   data["cmd"].(string),
		env:   toStringMapString(data["env"]),
	}
	return &c
}

// Type ...
func (c *Command) Type() DefinitionType {
	return CommandType
}

// Ident ...
func (c *Command) Ident() string {
	return c.ident
}

// Cmd ...
func (c *Command) Cmd() string {
	return c.cmd
}

// Env ...
func (c *Command) Env() map[string]string {
	return c.env
}

func toStringMapString(originalData interface{}) map[string]string {
	acc := make(map[string]string)
	for keyI, valueI := range originalData.(map[interface{}]interface{}) {
		acc[keyI.(string)] = valueI.(string)
	}
	return acc
}
