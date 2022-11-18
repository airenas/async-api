package miniofs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert.Nil(t, validate(Options{URL: "olia", User: "olia", Bucket: "olia"}))
	assert.NotNil(t, validate(Options{URL: "", User: "olia", Bucket: "olia"}))
	assert.NotNil(t, validate(Options{URL: "olia", User: "", Bucket: "olia"}))
	assert.NotNil(t, validate(Options{URL: "olia", User: "olia", Bucket: ""}))
}
