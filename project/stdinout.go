package project

type stdinout struct {
	defType DefinitionType
	ident   string
}

// Stdin ...
var Stdin stdinout

// Stdout ...
var Stdout stdinout

func init() {
	Stdin = stdinout{StdinType, "<stdin>"}
	Stdout = stdinout{StdoutType, "<stdout>"}
}

// Type ...
func (sio *stdinout) Type() DefinitionType {
	return sio.defType
}

// Ident ...
func (sio *stdinout) Ident() string {
	return sio.ident
}
