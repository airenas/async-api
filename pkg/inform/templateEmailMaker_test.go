package inform

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateInit_Fails(t *testing.T) {
	m, err := newTemplateEmailMaker(viper.New(), testTemplate())
	assert.NotNil(t, err, "Error expected")
	assert.Nil(t, m)
}

func TestTemplateInit_FailsTemplate(t *testing.T) {
	v := viper.New()
	v.Set("mail.url", "url")
	_, err := newTemplateEmailMaker(v, "{{ xxx}")
	assert.NotNil(t, err)
}

func TestTemplateInit_OK(t *testing.T) {
	v := viper.New()
	v.Set("mail.url", "url")
	v.Set("smtp.username", "email")
	m, err := newTemplateEmailMaker(v, testTemplate())
	assert.Nil(t, err)
	assert.Equal(t, "url", m.url)
}

func TestTemplateEmail(t *testing.T) {
	v := viper.New()
	v.Set("mail.url", "url/{{ID}}")
	v.Set("smtp.username", "Olia")
	m, _ := newTemplateEmailMaker(v, testTemplate())
	require.NotNil(t, m)
	data := testEmailData()

	e, err := m.Make(data)
	require.Nil(t, err)
	assert.Equal(t, "Opps guys", e.Subject)
	assert.Contains(t, e.To, "email")
	assert.Equal(t, "txt id url/id "+data.MsgTime.Format(dateFormat), string(e.Text))
	assert.Equal(t, "<html>id url/id "+data.MsgTime.Format(dateFormat)+"</html>", string(e.HTML))
}

func TestTemplateEmail_Fail(t *testing.T) {
	v := viper.New()
	v.Set("mail.url", "url/{{ID}}")
	v.Set("smtp.username", "Olia")
	m, _ := newTemplateEmailMaker(v, testTemplate())
	require.NotNil(t, m)

	_, err := m.Make(testEmailData())
	require.Nil(t, err)

	m, _ = newTemplateEmailMaker(v, testTemplHTML+testTemplText)
	_, err = m.Make(testEmailData())
	require.NotNil(t, err)
	m, _ = newTemplateEmailMaker(v, testTemplSubject+testTemplText)
	_, err = m.Make(testEmailData())
	require.NotNil(t, err)
	m, _ = newTemplateEmailMaker(v, testTemplHTML+testTemplSubject)
	_, err = m.Make(testEmailData())
	require.NotNil(t, err)
}

const (
	testTemplSubject = `{{define "mail.Failed.subject"}}Opps guys{{end}}`
	testTemplText    = `{{define "mail.Failed.text"}}txt {{.ID}} {{.URL}} {{.Date}}{{end}}`
	testTemplHTML    = `{{define "mail.Failed.html"}}<html>{{.ID}} {{.URL}} {{.Date}}</html>{{end}}`
)

func testTemplate() string {
	return testTemplSubject + testTemplText + testTemplHTML
}

func testEmailData() *Data {
	return &Data{
		Email:   "email",
		ID:      "id",
		MsgType: "Failed",
		MsgTime: time.Now(),
	}
}
