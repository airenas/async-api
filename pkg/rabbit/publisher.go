package rabbit

import (
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

//Publisher publish events to rabbit mq broker
type Publisher struct {
	ChannelProvider *ChannelProvider
}

//NewPublisher initializes rabbit publisher
func NewPublisher(provider *ChannelProvider) *Publisher {
	return &Publisher{ChannelProvider: provider}
}

//Publish publish the message
func (sender *Publisher) Publish(id string, topic string) error {
	realTopic := sender.ChannelProvider.QueueName(topic)
	goapp.Log.Info().Msgf("Publishing event %s(%s)", realTopic, id)

	err := sender.ChannelProvider.RunOnChannelWithRetry(func(ch *amqp.Channel) error {
		return ch.Publish(
			realTopic, // exchange
			"",
			false, // mandatory
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(id),
			})
	})
	if err != nil {
		return errors.Wrap(err, "can't publish event")
	}
	return nil
}
