package mq

import (
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
)

type ServerOptionFunc func(*serverOptions)

// ServerOptionWithTimeout 超时控制
func ServerOptionWithTimeout(duration time.Duration) ServerOptionFunc {
	return func(options *serverOptions) {
		options.timeout = duration
	}
}

// ServerOptionWithMiddleware 中间件
func ServerOptionWithMiddleware(ms ...middleware.Middleware) ServerOptionFunc {
	return func(options *serverOptions) {
		options.ms = ms
	}
}

// ServerOptionWithDecodeFunc 解码函数
func ServerOptionWithDecodeFunc(f DecodeFunc) ServerOptionFunc {
	return func(options *serverOptions) {
		options.decodeFunc = f
	}
}

// ServerOptionWithErrHandler 运行错误处理器
func ServerOptionWithErrHandler(handler ErrorHandler) ServerOptionFunc {
	return func(options *serverOptions) {
		options.errHandler = handler
	}
}

// ServerOptionWithErrHandleFunc 运行错误处理函数
func ServerOptionWithErrHandleFunc(f ErrorHandlerFunc) ServerOptionFunc {
	return func(options *serverOptions) {
		options.errHandler = f
	}
}

// ServerOptionWithServerName 服务名定义
func ServerOptionWithServerName(name string) ServerOptionFunc {
	return func(options *serverOptions) {
		options.serverName = name
	}
}

// ServerOptionWithConcurrencyNum 并发数
func ServerOptionWithConcurrencyNum(num int32) ServerOptionFunc {
	return func(options *serverOptions) {
		options.cap = num
	}
}

// ServerOptionWithScaleThreshold 阈值 任务队列长度超过 这个数字 则会增加 go routine
func ServerOptionWithScaleThreshold(num int32) ServerOptionFunc {
	return func(options *serverOptions) {
		options.scaleThreshold = num
	}
}
