package mongo

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// Sanitize mongo input
func Sanitize(s string) string {
	return strings.Trim(s, " $/^\\")
}

// SkipNoDocErr checks error and skips if error is mongo.ErrNoDocuments
func SkipNoDocErr(err error) error {
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

// NewCollection initiates new collection for mongo
func NewCollection(pr *SessionProvider, tName string) (*mongo.Collection, context.Context, func(), error) {
	session, err := pr.NewSession()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can't init new session")
	}
	res := session.Client().Database(pr.store).Collection(tName)
	ctx, cancel := mongoContext()
	return res, ctx, func() {
		session.EndSession(context.Background())
		cancel()
	}, nil
}
