package inform

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/jordan-wright/email"
)

const dateFormat = "2006-01-02 15:04:05"

// TemplateEmailMaker makes email from provided template file
type TemplateEmailMaker struct {
	url       string
	c         *viper.Viper
	templates *template.Template
}

func NewTemplateEmailMaker(c *viper.Viper) (*TemplateEmailMaker, error) {
	tFile, err := getStringNonNil(c, "mail.template")
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadFile(tFile)
	if err != nil {
		return nil, err
	}
	return newTemplateEmailMaker(c, string(bytes))
}

func newTemplateEmailMaker(c *viper.Viper, tmplStr string) (*TemplateEmailMaker, error) {
	r := TemplateEmailMaker{c: c}
	var err error
	if r.url, err = getStringNonNil(c, "mail.url"); err != nil {
		return nil, err
	}
	r.templates, err = template.New("mail").Parse(tmplStr)
	if err != nil {
		return nil, errors.Wrapf(err, "can't parse template")
	}
	return &r, nil
}

//Make prepares the email for ID
func (maker *TemplateEmailMaker) Make(data *Data) (*email.Email, error) {
	return maker.make(data, maker.c)
}

func (maker *TemplateEmailMaker) prepareData(data *Data) *emailData {
	return &emailData{
		ID:   data.ID,
		URL:  strings.Replace(maker.url, "{{ID}}", data.ID, -1),
		Date: data.MsgTime.Format(dateFormat),
	}
}

type emailData struct {
	ID, URL, Date string
}

func (maker *TemplateEmailMaker) make(data *Data, c *viper.Viper) (*email.Email, error) {
	r := email.NewEmail()
	eData := maker.prepareData(data)
	if sub, err := maker.executeTempl("mail."+data.MsgType+".subject", eData); err != nil {
		return nil, err
	} else {
		r.Subject = string(sub)
	}
	var err error
	if r.Text, err = maker.executeTempl("mail."+data.MsgType+".text", eData); err != nil {
		return nil, err
	}
	if r.HTML, err = maker.executeTempl("mail."+data.MsgType+".html", eData); err != nil {
		return nil, err
	}
	r.To = []string{data.Email}
	r.From, err = getStringNonNil(c, "smtp.username")
	return r, err
}

func (maker *TemplateEmailMaker) executeTempl(name string, ed *emailData) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := maker.templates.ExecuteTemplate(buf, name, ed); err != nil {
		return nil, errors.Wrapf(err, "temlate %s", name)
	}
	return buf.Bytes(), nil
}
