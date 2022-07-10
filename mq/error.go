package mq

import "github.com/go-kratos/kratos/v2/log"

type ErrorHandler interface {
	Handle(err error)
}

type ErrorHandlerFunc func(err error)

func (f ErrorHandlerFunc) Handle(err error) {
	f(err)
}

func DefaultErrorHandler(err error) {
	log.Error(err)
}
