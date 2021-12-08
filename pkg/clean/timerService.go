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

type timerServiceData struct {
	runEvery    time.Duration
	cleaner     Cleaner
	idsProvider OldIDsProvider
}

func StartCleanTimer(ctx context.Context, data *timerServiceData) (<-chan struct{}, error) {
	if data.runEvery < time.Minute {
		return nil, errors.Errorf("wrong run every duration %s, expected >= 1m", data.runEvery.String())
	}
	goapp.Log.Infof("Starting timer service every %v", data.runEvery)
	res := make(chan struct{}, 2)
	go func() {
		defer close(res)
		serviceLoop(ctx, data)
	}()
	return res, nil
}

func serviceLoop(ctx context.Context, data *timerServiceData) {
	ticker := time.NewTicker(data.runEvery)
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

func doClean(data *timerServiceData) {
	goapp.Log.Info("Running cleaning")
	ids, err := data.idsProvider.GetExpired()
	if err != nil {
		goapp.Log.Error(err)
	}
	goapp.Log.Infof("Got %d IDs to clean", len(ids))
	for _, id := range ids {
		err = data.cleaner.Clean(id)
		if err != nil {
			goapp.Log.Error(err)
		}
	}
}
