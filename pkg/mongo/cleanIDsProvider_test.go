package mongo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_getID(t *testing.T) {
	tests := []struct {
		name    string
		args    bson.M
		want    string
		wantErr bool
	}{
		{name: "Get", args: bson.M{"ID": "olia"}, want: "olia", wantErr: false},
		{name: "Fail", args: bson.M{"ID": 10}, want: "", wantErr: true},
		{name: "Fail", args: bson.M{"ID1": "tata"}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getID(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("getID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isOld(t *testing.T) {
	type args struct {
		m          bson.M
		expireDate time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "expired", args: args{m: bson.M{"ID": "olia"}, expireDate: time.Now().Add(time.Minute)}, want: false},
		{name: "expired", args: args{m: bson.M{"_id": primitive.NewObjectIDFromTimestamp(time.Now())},
			expireDate: time.Now().Add(time.Minute)}, want: true},
		{name: "expired", args: args{m: bson.M{"_id": primitive.NewObjectIDFromTimestamp(time.Now().Add(-time.Hour))},
			expireDate: time.Now().Add(-2 * time.Hour)}, want: false},
		{name: "expired", args: args{m: bson.M{"_id": primitive.NewObjectIDFromTimestamp(time.Now().Add(-3 * time.Hour))},
			expireDate: time.Now().Add(-2 * time.Hour)}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOld(tt.args.m, tt.args.expireDate); got != tt.want {
				t.Errorf("isOld() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCleanIDsProvider(t *testing.T) {
	type args struct {
		sessionProvider *SessionProvider
		expireDuration  time.Duration
		table           string
	}
	tests := []struct {
		name    string
		args    args
		want    *CleanIDsProvider
		wantErr bool
	}{
		{name: "OK", args: args{expireDuration: time.Hour, table: "table"}, wantErr: false},
		{name: "OK", args: args{expireDuration: time.Minute, table: "table"}, wantErr: false},
		{name: "Fail", args: args{expireDuration: time.Second, table: "table"}, wantErr: true},
		{name: "Fail", args: args{expireDuration: time.Hour, table: ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCleanIDsProvider(tt.args.sessionProvider, tt.args.expireDuration, tt.args.table)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCleanIDsProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}
