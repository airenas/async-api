package mocks

import (
	"testing"

	"github.com/petergtz/pegomock"
)

//go:generate pegomock generate --package=mocks --output=cleaner.go -m github.com/airenas/async-api/pkg/clean Cleaner

//go:generate pegomock generate --package=mocks --output=oldIDs.go -m github.com/airenas/async-api/pkg/clean OldIDsProvider


//AttachMockToTest register pegomock verification to be passed to testing engine
func AttachMockToTest(t *testing.T) {
	pegomock.RegisterMockFailHandler(handleByTest(t))
}

func handleByTest(t *testing.T) pegomock.FailHandler {
	return func(message string, callerSkip ...int) {
		if message != "" {
			t.Error(message)
		}
	}
}
