package actor

// Message ...
type Message struct {
	Type MessageType
	Data []byte
}

// NewDataMessage ...
func NewDataMessage(data []byte) Message {
	return Message{
		Type: DataMessage,
		Data: data,
	}
}

// NewCloseMessage ...
func NewCloseMessage() Message {
	return Message{
		Type: CloseInbox,
	}
}

// MessageType ...
type MessageType int

// MessageType ...
const (
	DataMessage MessageType = iota
	CloseInbox
)
