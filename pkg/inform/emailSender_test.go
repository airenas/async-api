package inform

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetType(t *testing.T) {
	testGetType(t, "", SMTP_PLAIN, false)
	testGetType(t, "olia", "", true)
	testGetType(t, SMTP_PLAIN, SMTP_PLAIN, false)
	testGetType(t, SMTP_LOGIN, SMTP_LOGIN, false)
	testGetType(t, SMTP_NOAUTH, SMTP_NOAUTH, false)
	testGetType(t, "no_auth", SMTP_NOAUTH, false)
}

func TestFullHost(t *testing.T) {
	se := SimpleEmailSender{}
	se.host = "net.olia"
	se.port = 445
	assert.Equal(t, "net.olia:445", se.getFullHost())
}

func testGetType(t *testing.T, v, exp string, expErr bool) {
	m, err := getType(v)
	assert.Equal(t, expErr, err != nil)
	assert.Equal(t, exp, m)
}

func TestNewSimpleEmailSender(t *testing.T) {
	type args struct {
		c *viper.Viper
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Init", args: args{c: newTestViper([][2]string{{"smtp.host", "host"}, {"smtp.port", "8000"}})}, wantErr: false},
		{name: "No host", args: args{c: newTestViper([][2]string{{"smtp.host", ""}, {"smtp.port", "8000"}})}, wantErr: true},
		{name: "No port", args: args{c: newTestViper([][2]string{{"smtp.host", "host"}, {"smtp.port", ""}})}, wantErr: true},
		{name: "Wrong type", args: args{c: newTestViper([][2]string{{"smtp.host", "host"}, {"smtp.port", "8000"},
			{"smtp.type", "PLAIN"}})}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSimpleEmailSender(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSimpleEmailSender() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func newTestViper(s [][2]string) *viper.Viper {
	res := viper.New()
	for _, c := range s {
		res.Set(c[0], c[1])
	}
	return res
}
