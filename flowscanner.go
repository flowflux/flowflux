package main

// FlowScanner ...
type FlowScanner interface {
	Scan() bool
	Message() []byte
	Err() error
}
