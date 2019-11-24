package actor

// Actable ...
type Actable interface {
	Next(next ...Actable)
	Receive(message Message)
	Start()
}
