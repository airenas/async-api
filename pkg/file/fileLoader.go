package file

import (
	"os"
	"path/filepath"

	"github.com/airenas/async-api/pkg/api"
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/pkg/errors"
)

//OpenFileReadFunc declares function to open file by name and return Reader
type OpenFileReadFunc func(fileName string) (api.FileRead, error)

// LocalLoader loads file on local disk
type LocalLoader struct {
	// StoragePath is the main folder to save into
	Path     string
	OpenFunc OpenFileReadFunc
}

//NewLocalLoader creates LocalLoader instance
func NewLocalLoader(path string) (*LocalLoader, error) {
	goapp.Log.Infof("Init Local File Loader at: %s", path)
	if path == "" {
		return nil, errors.New("no path provided")
	}
	f := LocalLoader{Path: path, OpenFunc: openFileForRead}
	return &f, nil
}

// Load loads file from disk
func (fs LocalLoader) Load(name string) (api.FileRead, error) {
	fileName := filepath.Join(fs.Path, name)
	f, err := fs.OpenFunc(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "can't open file %s", fileName)
	}
	return f, nil
}

func openFileForRead(fileName string) (api.FileRead, error) {
	return os.Open(fileName)
}
