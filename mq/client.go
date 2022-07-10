package mq

import (
	"context"
)

type ClientOptionFunc func(*clientOptions)

type clientOptions struct {
	encodeFunc EncodeFunc
}

type Client struct {
	producer Producer

	options clientOptions
}

func NewClient(producer Producer, opts ...ClientOptionFunc) *Client {
	options := &clientOptions{
		encodeFunc: DefaultEncodeFunc,
	}
	for _, f := range opts {
		f(options)
	}
	return &Client{producer: producer, options: *options}
}

func (x *Client) Invoke(ctx context.Context, topic string, args interface{}) error {
	message, err := x.options.encodeFunc(ctx, args)
	if err != nil {
		return err
	}
	return x.producer.Publish(ctx, topic, message)
}
