package miniofs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert.Nil(t, validate(Options{url: "olia", user: "olia", bucket: "olia"}))
	assert.NotNil(t, validate(Options{url: "", user: "olia", bucket: "olia"}))
	assert.NotNil(t, validate(Options{url: "olia", user: "", bucket: "olia"}))
	assert.NotNil(t, validate(Options{url: "olia", user: "olia", bucket: ""}))
}
