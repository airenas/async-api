package clean

import (
	"context"
	"time"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/pkg/errors"
)

//OldIDsProvider return old ids for cleaning service
type OldIDsProvider interface {
	GetExpired() ([]string, error)
}

type Cleaner interface {
	Clean(ID string) error
}

type TimerData struct {
	RunEvery    time.Duration
	Cleaner     Cleaner
	IDsProvider OldIDsProvider
}

func StartCleanTimer(ctx context.Context, data *TimerData) (<-chan struct{}, error) {
	if data.RunEvery < time.Minute {
		return nil, errors.Errorf("wrong run every duration %s, expected >= 1m", data.RunEvery.String())
	}
	goapp.Log.Infof("Starting timer service every %v", data.RunEvery)
	res := make(chan struct{}, 2)
	go func() {
		defer close(res)
		serviceLoop(ctx, data)
	}()
	return res, nil
}

func serviceLoop(ctx context.Context, data *TimerData) {
	ticker := time.NewTicker(data.RunEvery)
	// run on startup
	doClean(data)
	for {
		select {
		case <-ticker.C:
			doClean(data)
		case <-ctx.Done():
			ticker.Stop()
			goapp.Log.Infof("Stopped timer service")
			return
		}
	}
}

func doClean(data *TimerData) {
	goapp.Log.Info("Running cleaning")
	ids, err := data.IDsProvider.GetExpired()
	if err != nil {
		goapp.Log.Error(err)
	}
	goapp.Log.Infof("Got %d IDs to clean", len(ids))
	for _, id := range ids {
		err = data.Cleaner.Clean(id)
		if err != nil {
			goapp.Log.Error(err)
		}
	}
}
