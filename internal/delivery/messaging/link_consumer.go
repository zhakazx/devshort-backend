package messaging

import (
	"devshort-backend/internal/model"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type LinkConsumer struct {
	Log *logrus.Logger
}

func NewLinkConsumer(log *logrus.Logger) *LinkConsumer {
	return &LinkConsumer{
		Log: log,
	}
}

func (c LinkConsumer) Consume(message *sarama.ConsumerMessage) error {
	LinkEvent := new(model.LinkEvent)
	if err := json.Unmarshal(message.Value, LinkEvent); err != nil {
		c.Log.WithError(err).Error("error unmarshalling Link event")
		return err
	}

	// TODO process event
	c.Log.Infof("Received topic links with event: %v from partition %d", LinkEvent, message.Partition)
	return nil
}