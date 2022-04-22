package file

import (
	"io/fs"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewOldDirProvider(t *testing.T) {
	type args struct {
		expireDuration time.Duration
		dir            string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "OK", args: args{expireDuration: time.Hour, dir: "olia"}, wantErr: false},
		{name: "No dir", args: args{expireDuration: time.Hour, dir: ""}, wantErr: true},
		{name: "Wrong exp duration", args: args{expireDuration: time.Millisecond * 900, dir: "olia"}, wantErr: true},
		{name: "Wrong exp duration", args: args{expireDuration: -time.Minute * 5, dir: "olia"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOldDirProvider(tt.args.expireDuration, tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOldDirProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("NewOldDirProvider() no instance")
				return
			}
		})
	}
}

func Test_filterExpired(t *testing.T) {
	type args struct {
		before time.Time
		files  []fs.FileInfo
	}
	now := time.Now()
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "Empty", args: args{before: now.Add(-time.Minute * 10),
			files: []fs.FileInfo{newMockFile("olia", now.Add(-time.Minute*2))}}, want: nil},
		{name: "Filters", args: args{before: now.Add(-time.Minute * 10),
			files: []fs.FileInfo{newMockFile("olia", now.Add(-time.Minute*11))}}, want: []string{"olia"}},
		{name: "Several", args: args{before: now.Add(-time.Minute * 10),
			files: []fs.FileInfo{newMockFile("olia", now.Add(-time.Minute*11)),
				newMockFile("olia1", now.Add(-time.Minute*8)),
				newMockFile("olia2", now.Add(-time.Minute*9))}}, want: []string{"olia"}},
		{name: "Several returns", args: args{before: now.Add(-time.Minute * 5),
			files: []fs.FileInfo{newMockFile("olia", now.Add(-time.Minute*11)),
				newMockFile("olia1", now.Add(-time.Minute*8)),
				newMockFile("olia2", now.Add(-time.Minute*9))}}, want: []string{"olia", "olia1", "olia2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterExpired(tt.args.before, tt.args.files)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newMockFile(f string, t time.Time) fs.FileInfo {
	return mockFileInfo{file: f, mod: t}
}

type mockFileInfo struct {
	file string
	mod  time.Time
}

func (mfi mockFileInfo) Name() string       { return mfi.file }
func (mfi mockFileInfo) Size() int64        { return int64(8) }
func (mfi mockFileInfo) Mode() os.FileMode  { return os.ModePerm }
func (mfi mockFileInfo) ModTime() time.Time { return mfi.mod }
func (mfi mockFileInfo) IsDir() bool        { return true }
func (mfi mockFileInfo) Sys() interface{}   { return nil }
