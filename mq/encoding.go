package mq

import (
	"context"

	messagev1 "github.com/Ccheers/kratos-mq/mq/message/v1"
	"github.com/go-kratos/kratos/v2/encoding"
	ejson "github.com/go-kratos/kratos/v2/encoding/json"
	eproto "github.com/go-kratos/kratos/v2/encoding/proto"
	"github.com/go-kratos/kratos/v2/metadata"
	"google.golang.org/protobuf/proto"
)

const contentTypeKey = "Content-Type"

type EncodeFunc func(ctx context.Context, args interface{}) (Message, error)

type DecodeFunc func(ctx context.Context, message Message, args interface{}) error

func DefaultEncodeFunc(ctx context.Context, args interface{}) (Message, error) {
	var codec encoding.Codec
	if _, ok := args.(proto.Message); ok {
		codec = encoding.GetCodec(eproto.Name)
	} else {
		codec = encoding.GetCodec(ejson.Name)
	}
	md, ok := metadata.FromServerContext(ctx)
	if ok {
		codec = encoding.GetCodec(md.Get(contentTypeKey))
	} else {
		md = metadata.New()
	}
	md.Set(contentTypeKey, codec.Name())

	b, err := codec.Marshal(args)
	if err != nil {
		return nil, err
	}
	return messagev1.NewMessage(b, messagev1.WithMetadata(md)), nil
}

func DefaultDecodeFunc(ctx context.Context, message Message, args interface{}) error {
	md := message.Metadata()
	codec := encoding.GetCodec(md.Get(contentTypeKey))
	err := codec.Unmarshal(message.Payload(), args)
	if err != nil {
		return err
	}
	return nil
}
