package mq

import (
	"context"
	"time"

	"github.com/ccheers/xpkg/sync/routinepool"
	"github.com/go-kratos/kratos/v2/middleware"
)

type ServerOptionFunc func(*serverOptions)

type serverOptions struct {
	// 配置
	timeout time.Duration

	// 逻辑功能
	ms         []middleware.Middleware
	decodeFunc DecodeFunc
	errHandler ErrorHandler

	// routine pool config
	serverName     string
	cap            int32
	scaleThreshold int32 // 阈值 任务队列长度超过 这个数字 则会增加 go routine
}

type Server struct {
	consumer Consumer

	options serverOptions

	pool routinepool.Pool
}

func NewServer(consumer Consumer, opts ...ServerOptionFunc) *Server {
	options := &serverOptions{
		timeout:    time.Second * 3,
		decodeFunc: DefaultDecodeFunc,
		errHandler: ErrorHandlerFunc(DefaultErrorHandler),
	}
	for _, f := range opts {
		f(options)
	}
	return &Server{
		consumer: consumer,
		options:  *options,
		pool: routinepool.NewPool(options.serverName, options.cap, &routinepool.Config{
			ScaleThreshold: options.scaleThreshold,
		}),
	}
}

type HandleFunc func(ctx context.Context, message Message)

func (f HandleFunc) Handle(ctx context.Context, message Message) {
	f(ctx, message)
}

type Handler interface {
	Handle(ctx context.Context, message Message)
}

func (x *Server) Subscriber(ctx context.Context, topic string, channel string, handler Handler, ms ...middleware.Middleware) error {
	ch, err := x.consumer.Subscribe(ctx, topic, channel)
	if err != nil {
		return err
	}
	_ctx, cancel := context.WithTimeout(context.Background(), x.options.timeout)
	_ctx = MiddlewareWithContext(_ctx, append(x.options.ms, ms...)...)
	x.pool.CtxGo(_ctx, func(ctx context.Context) {
		defer cancel()
		var msg Message
		for {
			select {
			case <-ctx.Done():
				return
			case msg = <-ch:
				handler.Handle(ctx, msg)
			default:
				time.Sleep(time.Second)
			}
		}
	})
	return nil
}

func (x *Server) DecodeFunc() DecodeFunc {
	return x.options.decodeFunc
}

func (x *Server) ErrHandler() ErrorHandler {
	return x.options.errHandler
}

func (x *Server) Start(ctx context.Context) error {
	return nil
}

func (x *Server) Stop(ctx context.Context) error {
	return x.consumer.Close(ctx)
}
