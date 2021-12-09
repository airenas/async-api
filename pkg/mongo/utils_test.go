package mongo

import (
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestSanitize(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{name: "Keep", args: "IDS134", want: "IDS134"},
		{name: "Trim", args: " IDS134", want: "IDS134"},
		{name: "Trim", args: " $^ IDS134", want: "IDS134"},
		{name: "Trim", args: "IDS134 ^$/\\", want: "IDS134"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sanitize(tt.args); got != tt.want {
				t.Errorf("Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkipNoDocErr(t *testing.T) {
	tests := []struct {
		name    string
		args    error
		wantErr error
	}{
		{name: "Skip", args: nil, wantErr: nil},
		{name: "Skip", args: mongo.ErrNoDocuments, wantErr: nil},
		{name: "Skip", args: mongo.ErrClientDisconnected, wantErr: mongo.ErrClientDisconnected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SkipNoDocErr(tt.args); err != tt.wantErr {
				t.Errorf("SkipNoDocErr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
