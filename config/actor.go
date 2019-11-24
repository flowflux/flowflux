package config

// ActorType ...
type ActorType int

// ActorType ...
const (
	StdinType ActorType = iota
	CommandType
	ForkType
	MergeType
	StdoutType
)

// Actor ...
type Actor interface {
	Type() ActorType
	ID() string
}
