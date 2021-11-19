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

// LocalSaver saves file on local disk
type LocalSaver struct {
	// StoragePath is the main folder to save into
	StoragePath  string
	OpenFileFunc OpenFileFunc
}

//NewLocalSaver creates LocalSaver instance
func NewLocalSaver(storagePath string) (*LocalSaver, error) {
	goapp.Log.Infof("Init Local File Storage at: %s", storagePath)
	if storagePath == "" {
		return nil, errors.New("no storage path provided")
	}
	if err := checkCreateDir(storagePath); err != nil {
		return nil, errors.Wrapf(err, "can't create dir %s", storagePath)
	}
	f := LocalSaver{StoragePath: storagePath, OpenFileFunc: openFile}
	return &f, nil
}

func checkCreateDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		goapp.Log.Infof("Trying to create storage directory at: %s", dir)
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

// Save saves file to disk
func (fs LocalSaver) Save(name string, reader io.Reader) error {
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
	if err := checkCreateDir(dir); err != nil {
		return nil, errors.Wrapf(err, "can't create dir '%s'", dir)
	}
	return os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
}
