package api

import (
	"io"
	"os"
)

// FileRead interface
type FileRead interface {
	io.ReadCloser
	io.Seeker
	Stat() (os.FileInfo, error)
}