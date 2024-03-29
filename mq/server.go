package mq

import (
	"context"
	"time"

	"github.com/ccheers/xpkg/sync/routinepool"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
)

type serverOptions struct {
	// meta
	serverName string

	// 配置
	timeout time.Duration

	// 逻辑功能
	ms         []middleware.Middleware
	decodeFunc DecodeFunc
	errHandler ErrorHandler

	// routine pool config
	cap            int32
	scaleThreshold int32 // 阈值 任务队列长度超过 这个数字 则会增加 go routine
}

type Server struct {
	ctx        context.Context
	cancelFunc func()

	consumer Consumer

	options serverOptions

	pool routinepool.Pool
}

func NewServer(ctx context.Context, consumer Consumer, opts ...ServerOptionFunc) *Server {
	options := &serverOptions{
		timeout:        time.Second * 3,
		decodeFunc:     DefaultDecodeFunc,
		errHandler:     ErrorHandlerFunc(DefaultErrorHandler),
		cap:            4,
		scaleThreshold: 4,
		ms: []middleware.Middleware{
			recovery.Recovery(
				recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
					log.Errorf("panic recover err=%v req=%+v", err, req)
					return nil
				}),
			),
		},
	}
	for _, f := range opts {
		f(options)
	}
	pool := routinepool.NewPool(options.serverName, options.cap, &routinepool.Config{
		ScaleThreshold: options.scaleThreshold,
	})
	pool.SetPanicHandler(func(ctx context.Context, err error) {
		options.errHandler.Handle(err)
	})
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		ctx:        ctx,
		cancelFunc: cancel,
		consumer:   consumer,
		options:    *options,
		pool:       pool,
	}
}

type HandleFunc func(ctx context.Context, message Message)

func (f HandleFunc) Handle(ctx context.Context, message Message) {
	f(ctx, message)
}

type Handler interface {
	Handle(ctx context.Context, message Message)
}

func (x *Server) Subscriber(topic string, channel string, handler Handler, ms ...middleware.Middleware) error {
	ch, err := x.consumer.Subscribe(x.ctx, topic, channel)
	if err != nil {
		return err
	}
	for i := 0; i < int(x.options.cap); i++ {
		_ctx := MiddlewareWithContext(x.ctx, append(x.options.ms, ms...)...)
		x.pool.CtxGo(_ctx, func(ctx context.Context) {
			var msg Message
			var ok bool
			for {
				select {
				case <-ctx.Done():
					log.Warnf("routine exit context canceled topic=%s, channel=%s", topic, channel)
					return
				case msg, ok = <-ch:
					if !ok {
						log.Warnf("channel closed topic=%s, channel=%s", topic, channel)
						return
					}
					ctx, cancel := context.WithTimeout(ctx, x.options.timeout)
					handler.Handle(ctx, msg)
					cancel()
				default:
					time.Sleep(time.Second)
				}
			}
		})
	}
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
	x.cancelFunc()
	return x.consumer.Close(ctx)
}
