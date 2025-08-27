package messaging

import (
	"devshort-backend/internal/model"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type LinkProducer struct {
	Producer[*model.LinkEvent]
}

func NewLinkProducer(producer sarama.SyncProducer, log *logrus.Logger) *LinkProducer {
	return &LinkProducer{
		Producer: Producer[*model.LinkEvent]{
			Producer: producer,
			Topic:    "links",
			Log:      log,
		},
	}
}
