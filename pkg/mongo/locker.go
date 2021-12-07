package mongo

import (
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Locker acquires lock in db
type Locker struct {
	SessionProvider *SessionProvider
	table           string
}

//NewLocker creates Locker instance
func NewLocker(sessionProvider *SessionProvider, table string) (*Locker, error) {
	if table == "" {
		return nil, errors.New("no table")
	}
	f := Locker{SessionProvider: sessionProvider, table: table}
	return &f, nil
}

//Lock locks record for sending email
func (ss *Locker) Lock(id string, lockKey string) error {
	goapp.Log.Infof("Locking %s: %s", id, lockKey)

	c, ctx, cancel, err := NewCollection(ss.SessionProvider, ss.table)
	if err != nil {
		return err
	}
	defer cancel()

	// make sure we have the record
	err = SkipNoDocErr(c.FindOneAndUpdate(ctx, bson.M{
		"$and": []bson.M{{"ID": Sanitize(id)}, {"key": lockKey}}},
		bson.M{"$setOnInsert": bson.M{"status": 0}},
		options.FindOneAndUpdate().SetUpsert(true)).Err())
	if err != nil {
		return errors.Wrap(err, "can't insert email lock table")
	}

	return c.FindOneAndUpdate(ctx, bson.M{
		"$and": []bson.M{{"ID": Sanitize(id)}, {"key": lockKey}, {"status": 0}}},
		bson.M{"$set": bson.M{"status": 1}}, options.FindOneAndUpdate().SetUpsert(false)).Err()
}

//UnLock marks record with specific value
func (ss *Locker) UnLock(id string, lockKey string, value *int) error {
	goapp.Log.Infof("Unlocking table %s: %s", id, lockKey)

	c, ctx, cancel, err := NewCollection(ss.SessionProvider, ss.table)
	if err != nil {
		return err
	}
	defer cancel()

	return c.FindOneAndUpdate(ctx, bson.M{
		"$and": []bson.M{{"ID": Sanitize(id)}, {"key": lockKey}, {"status": 1}}},
		bson.M{"$set": bson.M{"status": *value}}, options.FindOneAndUpdate().SetUpsert(false)).Err()
}
