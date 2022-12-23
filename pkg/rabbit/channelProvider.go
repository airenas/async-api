package rabbit

import (
	"sync"
	"time"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/cenkalti/backoff/v4"
	"github.com/streadway/amqp"

	"github.com/pkg/errors"
)

//ChannelProvider provider amqp channel
type ChannelProvider struct {
	url  string
	conn *amqp.Connection
	ch   *amqp.Channel
	m    sync.Mutex // struct field mutex
}

type runOnChannelFunc func(*amqp.Channel) error

//NewChannelProvider initializes channel provider
func NewChannelProvider(url, user, pass string) (*ChannelProvider, error) {
	if url == "" {
		return nil, errors.New("no broker url set")
	}
	if user != "" && pass == "" {
		return nil, errors.New("no broker password set")
	}
	return &ChannelProvider{url: prepareURL(url, user, pass)}, nil
}

func prepareURL(url, user, pass string) string {
	res := "amqp://"
	if user != "" {
		res += user + ":" + pass + "@"
	}
	return res + url
}

//Channel return cached channel or tries to connect to rabbit broker
func (pr *ChannelProvider) Channel() (*amqp.Channel, error) {
	pr.m.Lock()
	defer pr.m.Unlock()

	if pr.ch != nil {
		return pr.ch, nil
	}
	conn, err := dial(pr.url)
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to rabbit broker")
	}
	ch, err := conn.Channel()
	if err != nil {
		defer conn.Close()
		return nil, errors.Wrap(err, "can't create channel")
	}
	pr.conn = conn
	pr.ch = ch
	return pr.ch, nil
}

//RunOnChannelWithRetry invokes method on channel with retry
func (pr *ChannelProvider) RunOnChannelWithRetry(f runOnChannelFunc) error {
	ch, err := pr.Channel()
	if err != nil {
		return errors.Wrap(err, "can't init channel")
	}
	err = f(ch)
	if err != nil {
		goapp.Log.Info().Msgf("retry opening channel")
		pr.Close()
		ch, err = pr.Channel()
		if err != nil {
			return errors.Wrap(err, "can't init channel")
		}
		err = f(ch)
	}
	return err
}

//Close finalizes ChannelProvider
func (pr *ChannelProvider) Close() {
	pr.m.Lock()
	defer pr.m.Unlock()

	if pr.ch != nil {
		_ = pr.ch.Close()
	}
	if pr.conn != nil {
		_ = pr.conn.Close()
	}
	pr.ch = nil
	pr.conn = nil
}

//QueueName return queue name for channel, may append prefix
func (pr *ChannelProvider) QueueName(name string) string {
	return name
}

// Healthy checks if rabbit channel is open
func (pr *ChannelProvider) Healthy() error {
	_, err := pr.Channel()
	if err != nil {
		return errors.Wrap(err, "Can't create channel")
	}
	return nil
}

func dial(url string) (*amqp.Connection, error) {
	var res *amqp.Connection
	op := func() error {
		var err error
		goapp.Log.Info().Msg("Dial " + goapp.HidePass(url))
		res, err = amqp.Dial(url)
		return err
	}
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 2 * time.Minute
	err := backoff.Retry(op, bo)
	if err == nil {
		goapp.Log.Info().Msg("Connected to " + goapp.HidePass(url))
	}
	return res, err
}
