package usage

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/airenas/async-api/internal/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestNewRestorer(t *testing.T) {
	got, err := NewRestorer("url", "")
	assert.Nil(t, err)
	assert.NotNil(t, got)
	_, err = NewRestorer("", "")
	assert.NotNil(t, err)
}

func TestWorker_Do(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/tt/restore/m:rid", r.RequestURI)
		b, _ := io.ReadAll(r.Body)
		assert.Equal(t, `{"error":"err"}`, string(b))
		rw.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	got, err := NewRestorer(srv.URL, "")
	assert.Nil(t, err)
	err = got.Do(test.Ctx(t), "id1", "tt:m:rid", "err")
	assert.Nil(t, err)
}

func TestWorker_AddHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Key some-secret-key", r.Header.Get("Authorization"))
		assert.Equal(t, "/tt/restore/m:rid", r.RequestURI)
		b, _ := io.ReadAll(r.Body)
		assert.Equal(t, `{"error":"err"}`, string(b))
		rw.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	got, err := NewRestorer(srv.URL, "some-secret-key")
	assert.Nil(t, err)
	err = got.Do(test.Ctx(t), "id1", "tt:m:rid", "err")
	assert.Nil(t, err)
}

func TestWorker_Skip_NoRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Fail(t, "unexpected call")
	}))
	defer srv.Close()
	got, err := NewRestorer(srv.URL, "")
	assert.Nil(t, err)
	err = got.Do(test.Ctx(t), "id1", "", "err")
	assert.Nil(t, err)
}

func TestWorker_Fail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()
	got, err := NewRestorer(srv.URL, "")
	assert.Nil(t, err)
	err = got.Do(test.Ctx(t), "id1", "tt:m:rid", "err")
	assert.Error(t, err)
}

func Test_parse(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		wantS   string
		wantR   string
		wantErr bool
	}{
		{name: "Parse", args: "tts:m:olia", wantS: "tts", wantR: "m:olia", wantErr: false},
		{name: "No manual", args: "tts::olia", wantS: "tts", wantR: ":olia", wantErr: false},
		{name: "Fails", args: "", wantS: "", wantR: "", wantErr: true},
		{name: "Fails", args: "tts:", wantS: "", wantR: "", wantErr: true},
		{name: "Fails", args: ":mma", wantS: "", wantR: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantS {
				t.Errorf("parse() got = %v, want %v", got, tt.wantS)
			}
			if got1 != tt.wantR {
				t.Errorf("parse() got1 = %v, want %v", got1, tt.wantR)
			}
		})
	}
}
