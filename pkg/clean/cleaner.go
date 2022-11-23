package clean

import (
	"strings"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/pkg/errors"
)

// CleanerGroup is a list of cleaners
type CleanerGroup struct {
	Jobs []Cleaner
}

// Clean runs all cleaners in the group
func (c *CleanerGroup) Clean(ID string) error {
	failed := 0
	for _, job := range c.Jobs {
		err := job.Clean(ID)
		if err != nil {
			goapp.Log.Error().Err(err).Send()
			failed++
		}
	}
	if failed == len(c.Jobs) {
		return errors.New("all delete tasks failed")
	}
	return nil
}

// NewFileCleaners creates file cleaners based on provided paths
func NewFileCleaners(fs string, patterns []string) ([]*LocalFile, error) {
	result := make([]*LocalFile, 0)
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p != "" {
			fc, err := NewLocalFile(fs, p)
			if err != nil {
				return nil, err
			}
			result = append(result, fc)
		}
	}
	return result, nil
}
