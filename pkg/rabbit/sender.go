package rabbit

import (
	"encoding/json"

	"github.com/airenas/async-api/pkg/messages"
	"github.com/airenas/go-app/pkg/goapp"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

//Sender performs messages sending using rabbit mq broker
type Sender struct {
	ChannelProvider *ChannelProvider
}

type initFunc func(*ChannelProvider) error

//NewSender initializes rabbit sender
func NewSender(provider *ChannelProvider) *Sender {
	return &Sender{ChannelProvider: provider}
}

//Send sends the message
func (sender *Sender) Send(message messages.Message, queue string, replyQueue string) error {
	return sender.SendWithCorr(message, queue, replyQueue, "")
}

//SendWithCorr sends the message with correlationID
func (sender *Sender) SendWithCorr(message messages.Message, queue string, replyQueue string, corrID string) error {
	realQueue := sender.ChannelProvider.QueueName(queue)
	goapp.Log.Debug().Msgf("Sending message to %s", realQueue)

	msgBytes, err := getBytes(message)
	if err != nil {
		return errors.Wrap(err, "can't marshal message")
	}

	err = sender.ChannelProvider.RunOnChannelWithRetry(func(ch *amqp.Channel) error {
		return ch.Publish(
			"", // exchange
			realQueue,
			false, // mandatory
			false,
			amqp.Publishing{
				DeliveryMode:  amqp.Persistent,
				ContentType:   "application/json",
				Body:          msgBytes,
				ReplyTo:       replyQueue,
				CorrelationId: corrID,
			})
	})
	if err != nil {
		return errors.Wrap(err, "Can't send message")
	}
	return nil
}

func getBytes(msg messages.Message) ([]byte, error) {
	res, err := json.Marshal(msg)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal message")
	}
	return res, nil
}
