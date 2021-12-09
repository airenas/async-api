package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCleanRecord(t *testing.T) {
	type args struct {
		sessionProvider *SessionProvider
		table           string
	}
	tests := []struct {
		name    string
		args    args
		want    *CleanRecord
		wantErr bool
	}{
		{name: "OK", args: args{table: "table"}, wantErr: false},
		{name: "Fail", args: args{table: ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCleanRecord(tt.args.sessionProvider, tt.args.table)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCleanRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}
