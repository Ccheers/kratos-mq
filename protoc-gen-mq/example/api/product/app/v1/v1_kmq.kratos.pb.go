// Code generated by github.com/Ccheers/kratos-mq/protoc-gen-mq. DO NOT EDIT.

package v1

import (
	context "context"
	mq "github.com/Ccheers/kratos-mq/mq"
	metadata "github.com/go-kratos/kratos/v2/metadata"
	middleware "github.com/go-kratos/kratos/v2/middleware"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the Ccheers/protoc-gen-mq package it is being compiled against.

// blog service is a blog demo
type BlogServiceMQServer interface {
	MQ_CreateArticle(context.Context, *Article) error
}

func RegisterBlogServiceMQServer(svr *mq.Server, srv BlogServiceMQServer) error {
	var err error
	err = svr.Subscriber("tp1", "ch3", CreateArticle_0_MQHandler(svr, srv))
	if err != nil {
		return err
	}
	err = svr.Subscriber("tp2", "ch2", CreateArticle_1_MQHandler(svr, srv))
	if err != nil {
		return err
	}
	return nil
}
func CreateArticle_0_MQHandler(svr *mq.Server, srv BlogServiceMQServer) mq.HandleFunc {
	return func(ctx context.Context, message mq.Message) {

		var in Article
		var err error

		err = svr.DecodeFunc()(ctx, message, &in)
		if err != nil {
			svr.ErrHandler().Handle(err)
			return
		}

		md := metadata.New(nil)

		for k, v := range message.Metadata() {
			md.Set(k, v)
		}

		newCtx := metadata.NewServerContext(ctx, md)

		ms := mq.GetMiddlewareFromContext(ctx)

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			err := srv.MQ_CreateArticle(ctx, req.(*Article))
			return nil, err
		}

		_, err = middleware.Chain(ms...)(handler)(newCtx, &in)
		if err != nil {
			svr.ErrHandler().Handle(err)
			return
		}

	}
}
func CreateArticle_1_MQHandler(svr *mq.Server, srv BlogServiceMQServer) mq.HandleFunc {
	return func(ctx context.Context, message mq.Message) {

		var in Article
		var err error

		err = svr.DecodeFunc()(ctx, message, &in)
		if err != nil {
			svr.ErrHandler().Handle(err)
			return
		}

		md := metadata.New(nil)

		for k, v := range message.Metadata() {
			md.Set(k, v)
		}

		newCtx := metadata.NewServerContext(ctx, md)

		ms := mq.GetMiddlewareFromContext(ctx)

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			err := srv.MQ_CreateArticle(ctx, req.(*Article))
			return nil, err
		}

		_, err = middleware.Chain(ms...)(handler)(newCtx, &in)
		if err != nil {
			svr.ErrHandler().Handle(err)
			return
		}

	}
}

type BlogServiceMQClient interface {
	CreateArticle_tp1(context.Context, *Article) error
	CreateArticle_tp2(context.Context, *Article) error
}
type BlogServiceMQClientImpl struct {
	cc *mq.Client
}

func NewBlogServiceMQClient(cc *mq.Client) BlogServiceMQClient {
	return &BlogServiceMQClientImpl{cc: cc}
}
func (x *BlogServiceMQClientImpl) CreateArticle_tp1(ctx context.Context, req *Article) error {
	return x.cc.Invoke(ctx, "tp1", req)
}
func (x *BlogServiceMQClientImpl) CreateArticle_tp2(ctx context.Context, req *Article) error {
	return x.cc.Invoke(ctx, "tp2", req)
}
