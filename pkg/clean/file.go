package clean

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/pkg/errors"
)

// LocalFile is a struct for local file cleaner
type LocalFile struct {
	storagePath string
	pattern     string
}

// NewLocalFile creates file cleaner
func NewLocalFile(storagePath string, pattern string) (*LocalFile, error) {
	goapp.Log.Info().Msgf("Init Local File Storage Clean at: %s/%s", storagePath, pattern)
	if pattern == "" {
		return nil, errors.New("no pattern provided")
	}
	if !strings.Contains(pattern, "{ID}") {
		return nil, errors.New("pattern does not contain {ID}")
	}
	sP := ""
	if !strings.HasPrefix(pattern, "/") {
		if storagePath == "" {
			return nil, errors.New("no storage path provided")
		}
		sP = storagePath
	}
	f := LocalFile{storagePath: sP, pattern: pattern}
	return &f, nil
}

// Clean removes files matching the pattern
func (fs *LocalFile) Clean(ID string) error {
	fp := fs.getPath(ID)
	goapp.Log.Info().Msgf("Removing %s", fp)
	return remove(fp)
}

func remove(fn string) error {
	files, err := filepath.Glob(fn)
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
		goapp.Log.Info().Msgf("Removed %s", file)
	}
	return nil
}

func (fs *LocalFile) getPath(ID string) string {
	res := strings.ReplaceAll(fs.pattern, "{ID}", ID)
	if fs.storagePath != "" {
		res = path.Join(fs.storagePath, res)
	}
	return res
}
