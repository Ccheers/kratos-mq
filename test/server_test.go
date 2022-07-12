package test

import (
	"context"
	"testing"
	"time"

	"github.com/Ccheers/kratos-mq/mq"
	"github.com/Ccheers/kratos-mq/mq_impl/nsq"
	"github.com/Ccheers/kratos-mq/mq_impl/nsq/config"
	v1 "github.com/Ccheers/kratos-mq/mq_impl/testdata/hello/v1"
	"github.com/go-kratos/kratos/v2/log"
)

var nsqCfg = &config.Config{
	NsqAddrs:          []string{"localhost:4150"},
	DialTimeout:       60,
	ReadTimeout:       60,
	WriteTimeout:      60,
	BatchSize:         10,
	HeartbeatInterval: 30,
}

type mockHelloMQServer struct {
}

func (x *mockHelloMQServer) MQ_HelloWorld(ctx context.Context, request *v1.HelloWorldRequest) error {
	log.Info(request.Msg)
	return nil
}

func TestNewServer(t *testing.T) {
	consumer, err := nsq.NewConsumer(nsqCfg, log.DefaultLogger)
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
			svr := mq.NewServer(tt.args.consumer, tt.args.opts...)
			err = v1.RegisterHelloMQServer(tt.args.ctx, svr, &mockHelloMQServer{})
			if err != nil && !tt.wantErr {
				t.Fatal(err)
			}
			clientTest(t)
			time.Sleep(time.Second * 5)
		})
	}
}

func clientTest(t *testing.T) {
	producer, err := nsq.NewProducer(nsqCfg, log.DefaultLogger)
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
