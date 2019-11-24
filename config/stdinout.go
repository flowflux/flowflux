package config

// Stdinout ...
type Stdinout struct {
	defType ActorType
	ident   string
}

// Stdin ...
var Stdin *Stdinout

// Stdout ...
var Stdout *Stdinout

// // Stderr ...
// var Stderr *Stdinout

func init() {
	Stdin = &Stdinout{StdinType, "<stdin>"}
	Stdout = &Stdinout{StdoutType, "<stdout>"}
	// Stderr = &Stdinout{StdoutType, "<stderr>"}
}

// Type ...
func (sio *Stdinout) Type() ActorType {
	return sio.defType
}

// ID ...
func (sio *Stdinout) ID() string {
	return sio.ident
}
