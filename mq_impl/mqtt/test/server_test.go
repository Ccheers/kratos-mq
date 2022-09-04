package test

import (
	"context"
	"testing"
	"time"

	v1 "github.com/Ccheers/kratos-mq/internal/testdata/hello/v1"
	"github.com/Ccheers/kratos-mq/mq"
	"github.com/Ccheers/kratos-mq/mq_impl/mqtt"
	mqtt1 "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
)

const EmqxCnBroker = "tcp://broker-cn.emqx.io:1883"

var cfg = mqtt1.NewClientOptions().
	AddBroker(EmqxCnBroker)

type mockHelloMQServer struct{}

func (x *mockHelloMQServer) MQ_HelloWorld(ctx context.Context, request *v1.HelloWorldRequest) error {
	log.Info("hello world", request.Msg)
	return nil
}

func TestNewServer(t *testing.T) {
	consumer, err := mqtt.NewConsumer(cfg, log.DefaultLogger)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		consumer mq.Consumer
		opts     []mq.ServerOptionFunc
		ctx      context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *mq.Server
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				consumer: consumer,
				ctx:      context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := mq.NewServer(tt.args.ctx, tt.args.consumer, tt.args.opts...)
			err := v1.RegisterHelloMQServer(svr, &mockHelloMQServer{})
			if err != nil && !tt.wantErr {
				t.Fatal(err)
			}
			clientTest(t)
			time.Sleep(time.Second * 5)
		})
	}
}

func clientTest(t *testing.T) {
	producer, err := mqtt.NewProducer(cfg)
	if err != nil {
		t.Error(err)
		return
	}
	cli := mq.NewClient(producer)
	err = v1.NewHelloMQClient(cli).HelloWorld_tp1(context.Background(), &v1.HelloWorldRequest{
		Msg: "123",
	})
	if err != nil {
		t.Error(err)
		return
	}
}
