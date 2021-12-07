package inform

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/airenas/async-api/pkg/inform/auth"
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/jordan-wright/email"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	//SMTP_NOAUTH login using no authentication
	SMTP_NOAUTH = "NO_AUTH"
	//SMTP_PLAIN login using PLAIN authentication for google
	SMTP_PLAIN = "PLAIN_AUTH"
	//SMTP_LOGIN login using no authentication for other
	SMTP_LOGIN = "LOGIN"
)

//SimpleEmailSender uses standard esmtp lib to send emails
type SimpleEmailSender struct {
	sendPool *email.Pool
	authType string
	host     string
	port     int
}

func NewSimpleEmailSender(c *viper.Viper) (*SimpleEmailSender, error) {
	r := SimpleEmailSender{}
	var err error
	r.authType, err = getType(c.GetString("smtp.type"))
	if err != nil {
		return nil, errors.Wrap(err, "can't init smtp authentication type")
	}
	r.host = c.GetString("smtp.host")
	if r.host == "" {
		return nil, errors.New("no smtp host")
	}
	r.port = c.GetInt("smtp.port")
	if r.port <= 0 {
		return nil, errors.New("no smtp port")
	}
	if r.authType != SMTP_NOAUTH {
		r.sendPool, err = email.NewPool(r.getFullHost(), 1, newAuth(r.authType, c))
		if err != nil {
			return nil, err
		}
	}
	goapp.Log.Infof("SMTP auth type: %s", r.authType)
	goapp.Log.Infof("SMTP server: %s", r.getFullHost())
	return &r, nil
}

func newAuth(authType string, c *viper.Viper) smtp.Auth {
	if authType == SMTP_LOGIN {
		goapp.Log.Infof("Using custom login auth")
		return auth.LoginAuth(c.GetString("smtp.username"), c.GetString("smtp.password"))
	}
	goapp.Log.Infof("Using plain login auth ")
	return smtp.PlainAuth("", c.GetString("smtp.username"), c.GetString("smtp.password"),
		c.GetString("smtp.host"))
}

//Send sends email
func (s *SimpleEmailSender) Send(email *email.Email) error {
	if s.authType == SMTP_NOAUTH {
		return email.Send(s.getFullHost(), nil)
	}
	return s.sendPool.Send(email, 10*time.Second)
}

func (s *SimpleEmailSender) getFullHost() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

func getType(s string) (string, error) {
	su := strings.TrimSpace(strings.ToUpper(s))
	if su == "" {
		return SMTP_PLAIN, nil
	}
	values := []string{SMTP_NOAUTH, SMTP_PLAIN, SMTP_LOGIN}
	for _, st := range values {
		if st == su {
			return su, nil
		}
	}
	return "", errors.Errorf("Unknown smtp type '%s'. Allowed values: %v", s, values)
}
