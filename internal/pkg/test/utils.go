package test

import (
	"context"
	"testing"
	"time"
)


// Ctx creates test context
func Ctx(t *testing.T) context.Context {
	t.Helper()
	ctx, cf := context.WithTimeout(context.Background(), time.Second*20)
	t.Cleanup(func() { cf() })
	return ctx
}

// To convert interface to object
func To[T interface{}](val interface{}) T {
	if val == nil {
		var res T
		return res
	}
	return val.(T)
}
