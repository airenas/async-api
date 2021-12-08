package clean

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileCleaners(t *testing.T) {
	f, err := NewFileCleaners("/path", []string{"path1{ID}"})
	assert.Nil(t, err)
	assert.NotNil(t, f)
}

func TestSeveralFileCleaners(t *testing.T) {
	f, err := NewFileCleaners("/path", []string{"path1{ID}", "{ID}.txt"})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, 2, len(f))
}

func TestNewFileCleanersPath(t *testing.T) {
	f, err := NewFileCleaners("/path", []string{"path1{ID}"})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, 1, len(f))
	assert.Equal(t, "/path", f[0].storagePath)
	assert.Equal(t, "path1{ID}", f[0].pattern)
}

func TestNewFileCleaners_Fail(t *testing.T) {
	_, err := NewFileCleaners("/path", []string{"path"})
	assert.NotNil(t, err)
}
