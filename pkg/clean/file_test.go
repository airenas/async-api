package clean

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLocalFile(t *testing.T) {
	type args struct {
		storagePath string
		pattern     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "init", args: args{storagePath: "path", pattern: "{ID}"}, wantErr: false},
		{name: "No ID", args: args{storagePath: "path", pattern: "{}"}, wantErr: true},
		{name: "No Path", args: args{storagePath: "", pattern: "ID"}, wantErr: true},
		{name: "No pattern", args: args{storagePath: "path", pattern: ""}, wantErr: true},
		{name: "No path, full pattern", args: args{storagePath: "", pattern: "/{ID}.txt"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLocalFile(tt.args.storagePath, tt.args.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLocalFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestLocalFile_getPath(t *testing.T) {
	type fields struct {
		storagePath string
		pattern     string
	}
	type args struct {
		ID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "Get", fields: fields{storagePath: "aa", pattern: "{ID}.txt"}, args: args{ID: "11"}, want: "aa/11.txt"},
		{name: "Get", fields: fields{storagePath: "/aa", pattern: "olia/{ID}.txt"}, args: args{ID: "11"}, want: "/aa/olia/11.txt"},
		{name: "No storage", fields: fields{storagePath: "", pattern: "/olia/{ID}.txt"}, args: args{ID: "11"}, want: "/olia/11.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &LocalFile{
				storagePath: tt.fields.storagePath,
				pattern:     tt.fields.pattern,
			}
			if got := fs.getPath(tt.args.ID); got != tt.want {
				t.Errorf("LocalFile.getPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
