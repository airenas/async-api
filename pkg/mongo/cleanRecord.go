package mongo

import (
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// CleanRecord deletes mongo table record
type CleanRecord struct {
	sessionProvider *SessionProvider
	table           string
}

// NewCleanRecord initiates object
func NewCleanRecord(sessionProvider *SessionProvider, table string) (*CleanRecord, error) {
	if table == "" {
		return nil, errors.New("no table")
	}
	if sessionProvider == nil {
		return nil, errors.New("no session provider")
	}
	f := CleanRecord{sessionProvider: sessionProvider, table: table}
	goapp.Log.Info().Msgf("Init Mongo table Clean for %s", table)
	return &f, nil
}

// Clean deletes record from table by ID
func (fs *CleanRecord) Clean(ID string) error {
	goapp.Log.Info().Msgf("Cleaning record for for %s[ID=%s]", fs.table, ID)

	c, ctx, cancel, err := NewCollection(fs.sessionProvider, fs.table)
	if err != nil {
		return err
	}
	defer cancel()

	info, err := c.DeleteMany(ctx, bson.M{"ID": ID})
	if err != nil {
		return errors.Wrap(err, "can't delete")
	}
	goapp.Log.Info().Msgf("Deleted %d", info.DeletedCount)
	return nil
}
