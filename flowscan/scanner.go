package flowscan

// Scanner ...
type Scanner interface {
	Scan() bool
	Message() []byte
	Err() error
}
