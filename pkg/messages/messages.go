package messages

import "time"

//Tag keeps key/value in message
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//Message base message interface for sending to queue
type Message interface {
}

//QueueMessage message going throuht broker
type QueueMessage struct {
	ID    string `json:"id"`
	Tags  []Tag  `json:"tags,omitempty"`
	Error string `json:"error,omitempty"`
}

const (
	// InformType_Started type when process started
	InformType_Started string = "Started"
	// InformType_Finished type when process finished
	InformType_Finished string = "Finished"
	// InformType_Failed type when process failed
	InformType_Failed string = "Failed"
)

//InformMessage message with inform information
type InformMessage struct {
	QueueMessage
	Type string    `json:"type"`
	At   time.Time `json:"at"`
}

//NewQueueMessageFromM copies message
func NewQueueMessageFromM(m *QueueMessage) *QueueMessage {
	return &QueueMessage{ID: m.ID, Tags: m.Tags}
}
