package mqtt

import (
	"context"
	"sync/atomic"

	"github.com/Ccheers/kratos-mq/mq"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ProducerImpl struct {
	cfg    *mqtt.ClientOptions
	client mqtt.Client

	status uint32
}

func NewProducer(cfg *mqtt.ClientOptions) (mq.Producer, error) {
	client := mqtt.NewClient(cfg)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &ProducerImpl{
		client: client,
		cfg:    cfg,
	}, nil
}

func (x *ProducerImpl) Publish(ctx context.Context, topic string, message mq.Message) error {
	b, err := message.Marshal()
	if err != nil {
		return err
	}
	const retained bool = false
	token := x.client.Publish(topic, x.cfg.WillQos, retained, b)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (x *ProducerImpl) Close(ctx context.Context) error {
	if !atomic.CompareAndSwapUint32(&x.status, statusRunning, statusClosed) {
		return nil
	}

	if !x.client.IsConnected() {
		return nil
	}
	x.client.Disconnect(0)
	return nil
}
