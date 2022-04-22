package file

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"time"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/pkg/errors"
)

// OldDirProvider returns old directories to remove from a system
type OldDirProvider struct {
	expireDuration time.Duration
	dir            string
}

// NewOldDirProvider creates OldDirProvider instances
func NewOldDirProvider(expireDuration time.Duration, dir string) (*OldDirProvider, error) {
	if expireDuration < time.Minute {
		return nil, errors.Errorf("wrong expireDuration %s, expected >= 1m", expireDuration.String())
	}
	if dir == "" {
		return nil, errors.New("no dir")
	}
	f := OldDirProvider{expireDuration: expireDuration, dir: dir}
	return &f, nil
}

// GetExpired return expired file nams
func (p *OldDirProvider) GetExpired() ([]string, error) {
	files, err := ioutil.ReadDir(p.dir)
	if err != nil {
		return nil, fmt.Errorf("can't read dir %s: %w", p.dir, err)
	}
	return filterExpired(time.Now().Add(-p.expireDuration), files), nil
}

func filterExpired(before time.Time, files []fs.FileInfo) []string {
	goapp.Log.Infof("Getting old files, time < %s", before.String())

	var res []string
	for _, f := range files {
		if f.ModTime().Before(before) {
			res = append(res, f.Name())
		}
	}
	return res
}
