package project

// Fork ...
type Fork struct {
	ident string
}

// NewFork ...
func NewFork(ident string, originalData interface{}) *Fork {
	f := Fork{
		ident: ident,
	}
	return &f
}

// Type ...
func (f *Fork) Type() DefinitionType {
	return ForkType
}

// Ident ...
func (f *Fork) Ident() string {
	return f.ident
}
