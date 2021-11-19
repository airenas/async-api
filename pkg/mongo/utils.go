package mongo

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func Sanitize(s string) string {
	return strings.Trim(s, " $/^\\")
}

func SkipNoDocErr(err error) error {
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

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
