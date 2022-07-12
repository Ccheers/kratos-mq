package nsq

import (
	"context"
	"math/rand"

	"github.com/Ccheers/kratos-mq/mq"
	"github.com/Ccheers/kratos-mq/mq_impl/nsq/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nsqio/go-nsq"
)

type ProducerImpl struct {
	producer *nsq.Producer
}

func NewProducer(c *config.Config, logger log.Logger) (mq.Producer, error) {
	// Instantiate a producer.
	cfg := NewNsqConfigWithConfig(c)

	producer, err := nsq.NewProducer(c.NsqAddrs[rand.Intn(len(c.NsqAddrs))], cfg)
	if err != nil {
		return nil, err
	}
	producer.SetLogger(newLogWrapper(logger), nsq.LogLevelInfo)

	return &ProducerImpl{
		producer: producer,
	}, nil
}

func (x *ProducerImpl) Publish(ctx context.Context, topic string, message mq.Message) error {
	b, err := message.Marshal()
	if err != nil {
		return err
	}
	return x.producer.Publish(topic, b)
}

func (x *ProducerImpl) Close(ctx context.Context) error {
	x.producer.Stop()
	return nil
}
