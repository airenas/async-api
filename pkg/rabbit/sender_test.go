package rabbit

import (
	"testing"

	"github.com/airenas/async-api/pkg/messages"
	"github.com/stretchr/testify/assert"
)

func TestGetBytes_Simple(t *testing.T) {
	m := messages.QueueMessage{ID: "id", Error: "err"}
	b, err := getBytes(m)
	assert.Nil(t, err)
	assert.Equal(t, "{\"id\":\"id\",\"recognizer\":\"rec\"}", string(b))
}

// func TestGetBytes_ResultMsg(t *testing.T) {
// 	m := messages.ResultMessage{QueueMessage: *messages.NewQueueMessage("id", "rec", nil), Result: "res"}
// 	b, err := getBytes(m)
// 	assert.Nil(t, err)
// 	assert.Equal(t, "{\"id\":\"id\",\"recognizer\":\"rec\",\"result\":\"res\"}", string(b))
// }

func TestGetBytes_Bytes(t *testing.T) {
	b, err := getBytes([]byte("olia"))
	assert.Nil(t, err)
	assert.Equal(t, "olia", string(b))
}

func TestGetBytes_String(t *testing.T) {
	b, err := getBytes("olia")
	assert.Nil(t, err)
	assert.Equal(t, "\"olia\"", string(b))
}
