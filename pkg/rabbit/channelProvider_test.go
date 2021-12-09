package rabbit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyQueueName(t *testing.T) {
	var prv ChannelProvider
	assert.Equal(t, "", prv.QueueName(""))
}

func TestNoPrefix(t *testing.T) {
	var prv ChannelProvider
	assert.Equal(t, "olia", prv.QueueName("olia"))
}

func TestNewChannelProvider(t *testing.T) {
	type args struct {
		url  string
		user string
		pass string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "OK", args: args{url: "url", user: "user", pass: "pass"}, wantErr: false},
		{name: "OK", args: args{url: "url", user: "", pass: ""}, wantErr: false},
		{name: "Fail", args: args{url: "", user: "user", pass: "pass"}, wantErr: true},
		{name: "Fail", args: args{url: "url", user: "user", pass: ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChannelProvider(tt.args.url, tt.args.user, tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChannelProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func Test_prepareURL(t *testing.T) {
	type args struct {
		url  string
		user string
		pass string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Full", args: args{url: "url", user: "user", pass: "pass"}, want: "amqp://user:pass@url"},
		{name: "No user", args: args{url: "url", user: "", pass: ""}, want: "amqp://url"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prepareURL(tt.args.url, tt.args.user, tt.args.pass); got != tt.want {
				t.Errorf("prepareURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
