package rabbit

import (
	"testing"

	"github.com/airenas/async-api/pkg/messages"
	"github.com/stretchr/testify/assert"
)

func TestGetBytes_Simple(t *testing.T) {
	m := messages.QueueMessage{ID: "id", Error: "err"}
	b, err := getBytes(&m)
	assert.Nil(t, err)
	assert.Equal(t, "{\"id\":\"id\",\"error\":\"err\"}", string(b))
}
