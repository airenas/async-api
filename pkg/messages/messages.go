package messages

//Tag keeps key/value in message
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//QueueMessage message going throuht broker
type QueueMessage struct {
	ID    string `json:"id"`
	Tags  []Tag  `json:"tags,omitempty"`
	Error string `json:"error,omitempty"`
}
