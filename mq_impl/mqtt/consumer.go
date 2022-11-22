package mqtt

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/Ccheers/kratos-mq/mq"
	"github.com/go-kratos/kratos/v2/log"
)

var _ mq.Consumer = (*ConsumerImpl)(nil)

type ConsumerImpl struct {
	cfg    *mqtt.ClientOptions
	client mqtt.Client
	logger *log.Helper

	mu           sync.Mutex
	consumerChan map[string]chan mq.Message

	status uint32
}

func NewConsumer(cfg *mqtt.ClientOptions, logger log.Logger) (mq.Consumer, error) {
	client := mqtt.NewClient(cfg)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &ConsumerImpl{
		client:       client,
		cfg:          cfg,
		logger:       log.NewHelper(log.With(logger, "module", "mqtt.consumer")),
		consumerChan: make(map[string]chan mq.Message),
		status:       statusRunning,
	}, nil
}

func (x *ConsumerImpl) Subscribe(ctx context.Context, topic string, channel string) (<-chan mq.Message, error) {
	x.mu.Lock()
	defer x.mu.Unlock()

	uniKey := x.generateKey(topic, channel)
	if ch, ok := x.consumerChan[uniKey]; ok {
		return ch, nil
	}

	ch := make(chan mq.Message, 1)
	token := x.client.Subscribe(topic, x.cfg.WillQos, func(client mqtt.Client, message mqtt.Message) {
		if atomic.LoadUint32(&x.status) == statusClosed {
			return
		}
		msg, err := mq.NewMessageFromByte(message.Payload())
		if err != nil {
			x.logger.Errorw("topic", topic, "channel", channel, "payload", string(message.Payload()), "err", err)
			return
		}
		ch <- msg
		message.Ack()
	})
	if token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	x.consumerChan[uniKey] = ch
	return ch, nil
}

func (x *ConsumerImpl) Close(ctx context.Context) error {
	if !atomic.CompareAndSwapUint32(&x.status, statusRunning, statusClosed) {
		return nil
	}

	x.mu.Lock()
	defer x.mu.Unlock()

	if !x.client.IsConnected() {
		return nil
	}
	x.client.Disconnect(0)

	for uniKey := range x.consumerChan {
		close(x.consumerChan[uniKey])
	}
	return nil
}

func (x *ConsumerImpl) generateKey(topic string, channel string) string {
	return fmt.Sprintf("%s::%s", topic, channel)
}
