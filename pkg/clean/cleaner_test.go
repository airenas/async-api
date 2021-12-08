package clean

import (
	"errors"
	"testing"

	"github.com/airenas/async-api/internal/pkg/test/mocks"
	"github.com/petergtz/pegomock"
	"github.com/stretchr/testify/assert"
)

func TestNewFileCleaners(t *testing.T) {
	f, err := NewFileCleaners("/path", []string{"path1{ID}"})
	assert.Nil(t, err)
	assert.NotNil(t, f)
}

func TestSeveralFileCleaners(t *testing.T) {
	f, err := NewFileCleaners("/path", []string{"path1{ID}", "{ID}.txt"})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, 2, len(f))
}

func TestNewFileCleanersPath(t *testing.T) {
	f, err := NewFileCleaners("/path", []string{"path1{ID}"})
	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, 1, len(f))
	assert.Equal(t, "/path", f[0].storagePath)
	assert.Equal(t, "path1{ID}", f[0].pattern)
}

func TestNewFileCleaners_Fail(t *testing.T) {
	_, err := NewFileCleaners("/path", []string{"path"})
	assert.NotNil(t, err)
}

func TestCleanerGroup_Clean(t *testing.T) {
	type fields struct {
		Jobs []Cleaner
	}
	type args struct {
		ID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "Fail", fields: fields{Jobs: []Cleaner{newCleanMock(true)}}, args: args{ID: "1"}, wantErr: true},
		{name: "OK", fields: fields{Jobs: []Cleaner{newCleanMock(false)}}, args: args{ID: "1"}, wantErr: false},
		{name: "Several OK", fields: fields{Jobs: []Cleaner{newCleanMock(false), newCleanMock(false)}}, args: args{ID: "1"}, wantErr: false},
		{name: "Some fail", fields: fields{Jobs: []Cleaner{newCleanMock(false), newCleanMock(true)}}, args: args{ID: "1"}, wantErr: false},
		{name: "All fail", fields: fields{Jobs: []Cleaner{newCleanMock(true), newCleanMock(true)}}, args: args{ID: "1"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CleanerGroup{
				Jobs: tt.fields.Jobs,
			}
			if err := c.Clean(tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("CleanerGroup.Clean() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func newCleanMock(fail bool) *mocks.MockCleaner {
	res := mocks.NewMockCleaner()
	var err error
	if fail {
		err = errors.New("olia")
	}
	pegomock.When(res.Clean(pegomock.AnyString())).ThenReturn(err)
	return res
}
