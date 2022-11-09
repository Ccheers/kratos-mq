package test

import (
	"context"
	"testing"
	"time"

	v1 "github.com/Ccheers/kratos-mq/internal/testdata/hello/v1"
	"github.com/Ccheers/kratos-mq/mq"
	"github.com/Ccheers/kratos-mq/mq_impl/diskq"
	"github.com/Ccheers/kratos-mq/mq_impl/diskq/config"
	"github.com/go-kratos/kratos/v2/log"
)

var nsqCfg = &config.Config{
	DataPath:        "/Users/eric/Downloads/1/",
	MaxBytesPerFile: 1024 * 1024,
	MinMsgSize:      0,
	MaxMsgSize:      1024,
	SyncEvery:       0,
	SyncTimeout:     1,
}

type mockHelloMQServer struct{}

func (x *mockHelloMQServer) MQ_HelloWorld(ctx context.Context, request *v1.HelloWorldRequest) error {
	log.Info(request.Msg)
	return nil
}

func TestNewServer(t *testing.T) {
	consumer, err := diskq.NewConsumer(nsqCfg, log.DefaultLogger)
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
			err = v1.RegisterHelloMQServer(svr, &mockHelloMQServer{})
			if err != nil && !tt.wantErr {
				t.Fatal(err)
			}
			clientTest(t)
			time.Sleep(time.Second * 5)
			svr.Stop(tt.args.ctx)
			time.Sleep(time.Second)
		})
	}
}

func clientTest(t *testing.T) {
	producer, err := diskq.NewProducer(nsqCfg, log.DefaultLogger)
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
