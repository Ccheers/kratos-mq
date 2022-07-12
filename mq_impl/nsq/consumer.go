package nsq

import (
	"context"
	"fmt"
	"sync"

	"github.com/Ccheers/kratos-mq/mq"
	"github.com/Ccheers/kratos-mq/mq_impl/nsq/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nsqio/go-nsq"
)

var _ mq.Consumer = (*ConsumerImpl)(nil)

type ConsumerImpl struct {
	nsqAddrs []string
	cfg      *nsq.Config

	logger *logWrapper

	mu           sync.Mutex
	consumerMap  map[string]*nsq.Consumer
	consumerChan map[string]chan mq.Message
}

func NewConsumer(c *config.Config, logger log.Logger) (mq.Consumer, error) {
	// Instantiate a consumer that will subscribe to the provided channel.
	cfg := NewNsqConfigWithConfig(c)
	return &ConsumerImpl{
		cfg:          cfg,
		nsqAddrs:     c.NsqAddrs,
		logger:       newLogWrapper(logger),
		consumerMap:  make(map[string]*nsq.Consumer),
		consumerChan: make(map[string]chan mq.Message),
	}, nil
}

func (x *ConsumerImpl) Subscribe(ctx context.Context, topic string, channel string) (<-chan mq.Message, error) {
	x.mu.Lock()
	defer x.mu.Unlock()

	uniKey := x.generateKey(topic, channel)
	if _, ok := x.consumerMap[uniKey]; ok {
		return nil, nil
	}

	consumer, err := nsq.NewConsumer(topic, channel, x.cfg)
	if err != nil {
		return nil, err
	}
	consumer.SetLoggerForLevel(x.logger, nsq.LogLevelInfo)

	ch := make(chan mq.Message, 1)
	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		msg, err := mq.NewMessageFromByte(message.Body)
		if err != nil {
			return err
		}
		ch <- msg
		message.Finish()
		return nil
	}))

	err = consumer.ConnectToNSQDs(x.nsqAddrs)
	if err != nil {
		return nil, err
	}

	x.consumerMap[uniKey] = consumer
	x.consumerChan[uniKey] = ch
	return ch, nil
}

func (x *ConsumerImpl) Close(ctx context.Context) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	for uniKey, consumer := range x.consumerMap {
		consumer.Stop()
		close(x.consumerChan[uniKey])
	}
	return nil
}

func (x *ConsumerImpl) generateKey(topic string, channel string) string {
	return fmt.Sprintf("%s::%s", topic, channel)
}
