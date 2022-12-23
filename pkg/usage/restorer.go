package usage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
)

// Restorer implements usage restore functionality
type Restorer struct {
	serviceURL string
	httpClient http.Client
}

// NewRestorer creates new restorer implementation
func NewRestorer(url string) (*Restorer, error) {
	if url == "" {
		return nil, errors.Errorf("no service URL")
	}
	res := &Restorer{serviceURL: url}
	res.httpClient = http.Client{Transport: &http.Transport{
		MaxIdleConns:        2,
		MaxIdleConnsPerHost: 2,
		IdleConnTimeout:     90 * time.Second,
		MaxConnsPerHost:     5,
	}}

	goapp.Log.Info().Str("URL", res.serviceURL).Msg("doorman-admin info")
	return res, nil
}

// Do tries to restore usage
func (w *Restorer) Do(ctx context.Context, msgID, reqID, errStr string) error {
	goapp.Log.Info().Str("ID", msgID).Str("requestID", reqID).Msg("doing usage restoratioon")
	if reqID == "" {
		goapp.Log.Warn().Msg("no requestID")
		return nil
	}
	service, rID, err := parse(reqID)
	if err != nil {
		return fmt.Errorf("wrong requestID format '%s': %w", reqID, err)
	}
	return w.invoke(ctx, service, rID, errStr)
}

func parse(s string) (string, string, error) {
	strs := strings.SplitN(s, ":", 2)
	if len(strs) != 2 || strs[0] == "" || strs[1] == "" {
		return "", "", fmt.Errorf("wrong format, expected 'srv:manual:requestID'")
	}
	return strs[0], strs[1], nil
}

type request struct {
	Error string `json:"error,omitempty"`
}

func (w *Restorer) invoke(ctx context.Context, service, requestID, errorMsg string) error {
	inp := request{Error: errorMsg}
	b, err := json.Marshal(inp)
	if err != nil {
		return err
	}

	_, err = goapp.InvokeWithBackoff(ctx, func() (any, bool, error) {
		req, err := http.NewRequest(http.MethodPost,
			fmt.Sprintf("%s/%s/restore/%s", w.serviceURL, service, requestID), bytes.NewReader(b))
		if err != nil {
			return nil, false, err
		}
		req.Header.Set("Content-Type", "application/json")
		
		ctx, cancelF := context.WithTimeout(ctx, time.Second*10)
		defer cancelF()
		req = req.WithContext(ctx)

		goapp.Log.Info().Str("URL", req.URL.String()).Msg("call")
		resp, err := w.httpClient.Do(req)
		if err != nil {
			return nil, goapp.IsRetryableErr(err), fmt.Errorf("can't invoke: %w", err)
		}
		defer func() {
			_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 10000))
			_ = resp.Body.Close()
		}()
		if err := goapp.ValidateHTTPResp(resp, 100); err != nil {
			err = fmt.Errorf("can't invoke '%s': %w", req.URL.String(), err)
			return nil, goapp.IsRetryableCode(resp.StatusCode), err
		}
		return nil, false, nil
	}, newSimpleBackoff())
	return err
}

func newSimpleBackoff() backoff.BackOff {
	res := backoff.NewExponentialBackOff()
	return backoff.WithMaxRetries(res, 3)
}
