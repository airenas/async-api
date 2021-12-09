package file

import (
	"errors"
	"testing"

	"github.com/airenas/async-api/internal/pkg/test/mocks"
	"github.com/airenas/async-api/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestLoads(t *testing.T) {
	fakeFile := fakeFile("content")
	fileLoader := LocalLoader{Path: "/data/",
		OpenFunc: func(file string) (api.FileRead, error) {
			return fakeFile, nil
		}}
	f, err := fileLoader.Load("file")
	assert.Nil(t, err)
	assert.NotNil(t, f)
}

func TestLoaderFailsOnNoOpen(t *testing.T) {
	fileLoader := LocalLoader{Path: "",
		OpenFunc: func(file string) (api.FileRead, error) {
			return nil, errors.New("olia")
		}}
	_, err := fileLoader.Load("file")
	assert.NotNil(t, err)
}

func TestLoaderChecksDirOnInit(t *testing.T) {
	_, err := NewLocalLoader("./")
	assert.Nil(t, err)
	_, err = NewLocalLoader("")
	assert.NotNil(t, err)
}

func fakeFile(c string) api.FileRead {
	return mocks.NewMockFileRead()
}
