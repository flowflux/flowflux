package config

// Fork ...
type Fork struct {
	Ident string
	Count int
}

// Type ...
func (f *Fork) Type() ActorType {
	return ForkType
}

// ID ...
func (f *Fork) ID() string {
	return f.Ident
}
