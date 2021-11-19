package file

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/pkg/errors"
)

//WriterCloser keeps Writer interface and close function
type WriterCloser interface {
	io.Writer
	Close() error
}

//OpenFileFunc declares function to open file by name and return Writer
type OpenFileFunc func(fileName string) (WriterCloser, error)

// LocalFileSaver saves file on local disk
type LocalFileSaver struct {
	// StoragePath is the main folder to save into
	StoragePath  string
	OpenFileFunc OpenFileFunc
}

//NewLocalFileSaver creates LocalFileSaver instance
func NewLocalFileSaver(storagePath string) (*LocalFileSaver, error) {
	goapp.Log.Infof("Init Local File Storage at: %s", storagePath)
	if storagePath == "" {
		return nil, errors.New("No storage path provided")
	}
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		goapp.Log.Infof("Trying to create storage directory at: %s", storagePath)
		err = os.MkdirAll(storagePath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	f := LocalFileSaver{StoragePath: storagePath, OpenFileFunc: openFile}
	return &f, nil
}

// Save saves file to disk
func (fs LocalFileSaver) Save(name string, reader io.Reader) error {
	if strings.Contains(name, "..") {
		return errors.New("wrong path " + name)
	}
	fileName := filepath.Join(fs.StoragePath, name)
	f, err := fs.OpenFileFunc(fileName)
	if err != nil {
		return errors.Wrapf(err, "can not create file %s", fileName)
	}
	defer f.Close()
	savedBytes, err := io.Copy(f, reader)
	if err != nil {
		return errors.Wrapf(err, "can not save file %s", fileName)
	}
	goapp.Log.Infof("Saved file %s. Size = %s b", fileName, strconv.FormatInt(savedBytes, 10))
	return nil
}

func openFile(fileName string) (WriterCloser, error) {
	dir := filepath.Dir(fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		goapp.Log.Infof("Trying to create storage directory at: %s", dir)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, errors.Wrapf(err, "can't create dir '%s'", dir)
		}
	}
	return os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
}