package clean

import (
	"context"
	"time"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/pkg/errors"
)

// OldIDsProvider return old ids for cleaning service
type OldIDsProvider interface {
	GetExpired(ctx context.Context) ([]string, error)
}

// Cleaner interface for one Clean job
type Cleaner interface {
	Clean(ctx context.Context, ID string) error
}

// TimerData keeps clean timer info
type TimerData struct {
	RunEvery    time.Duration
	Cleaner     Cleaner
	IDsProvider OldIDsProvider
}

// StartCleanTimer starts timer in loop for doing clean tasks
func StartCleanTimer(ctx context.Context, data *TimerData) (<-chan struct{}, error) {
	if data.RunEvery < time.Minute {
		return nil, errors.Errorf("wrong run every duration %s, expected >= 1m", data.RunEvery.String())
	}
	if data.Cleaner == nil {
		return nil, errors.Errorf("no cleaner")
	}
	if data.IDsProvider == nil {
		return nil, errors.Errorf("no IDs provider")
	}

	return startLoop(ctx, data), nil
}

func startLoop(ctx context.Context, data *TimerData) <-chan struct{} {
	goapp.Log.Info().Msgf("Starting timer service every %v", data.RunEvery)
	res := make(chan struct{}, 2)
	go func() {
		defer close(res)
		serviceLoop(ctx, data)
	}()
	return res
}

func serviceLoop(ctx context.Context, data *TimerData) {
	ticker := time.NewTicker(data.RunEvery)
	// run on startup
	doClean(ctx, data)
	for {
		select {
		case <-ticker.C:
			ctxInt, cf := context.WithTimeout(ctx, time.Second*60)
			doClean(ctxInt, data)
			cf()
		case <-ctx.Done():
			ticker.Stop()
			goapp.Log.Info().Msgf("Stopped timer service")
			return
		}
	}
}

func doClean(ctx context.Context, data *TimerData) {
	goapp.Log.Info().Msg("Running cleaning")
	ids, err := data.IDsProvider.GetExpired(ctx)
	if err != nil {
		goapp.Log.Error().Err(err).Send()
	}
	goapp.Log.Info().Int("count", len(ids)).Msg("Got IDs to clean")
	for _, id := range ids {
		err = data.Cleaner.Clean(ctx, id)
		if err != nil {
			goapp.Log.Error().Err(err).Send()
		}
	}
}
