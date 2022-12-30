package messages

import "time"

// Tag keeps key/value in message
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Message base message interface for sending to queue
type Message interface {
	GetID() string
}

// QueueMessage message going throuht broker
type QueueMessage struct {
	ID    string `json:"id"`
	Tags  []Tag  `json:"tags,omitempty"`
	Error string `json:"error,omitempty"`
}

const (
	// InformTypeStarted type when process started
	InformTypeStarted string = "Started"
	// InformTypeFinished type when process finished
	InformTypeFinished string = "Finished"
	// InformTypeFailed type when process failed
	InformTypeFailed string = "Failed"
)

// InformMessage message with inform information
type InformMessage struct {
	QueueMessage
	Type string    `json:"type"`
	At   time.Time `json:"at"`
}

// NewQueueMessageFromM copies message
func NewQueueMessageFromM(m *QueueMessage) *QueueMessage {
	return &QueueMessage{ID: m.ID, Tags: m.Tags}
}

// GetID returm message's ID
func (m *QueueMessage) GetID() string {
	return m.ID
}
