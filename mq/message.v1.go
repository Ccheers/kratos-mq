package mq

import (
	"errors"
	"fmt"
	"hash/crc32"
	"time"

	messagev1 "github.com/Ccheers/kratos-mq/mq/message/v1"
	"github.com/go-kratos/kratos/v2/metadata"
	"google.golang.org/protobuf/proto"
)

type MessageV1 messagev1.Message

var (
	ErrMessageIsNil     = errors.New("message is nil show alloc memory")
	ErrMessageIsInvalid = errors.New("message is invalid")
)

type MessageOption func(x *MessageV1)

func MessageOptionWithMetadata(md metadata.Metadata) MessageOption {
	return func(x *MessageV1) {
		x.Md = md
	}
}

var _ Message = (*MessageV1)(nil)

func NewMessage(payload Payload, opts ...MessageOption) Message {
	m := &MessageV1{Data: payload}
	// 校验和计算
	defer m.assignValidSum()
	for _, opt := range opts {
		opt(m)
	}
	// 兜底操作
	if m.Key == "" {
		m.Key = m.generateUniKey()
	}
	return m
}

func NewMessageFromByte(b []byte) (Message, error) {
	m := &MessageV1{}
	err := m.UnMarshal(b)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (x *MessageV1) Metadata() metadata.Metadata {
	return x.Md
}

func (x *MessageV1) Payload() Payload {
	return x.Data
}

func (x *MessageV1) Err() error {
	return fmt.Errorf(x.Error)
}

func (x *MessageV1) UniKey() string {
	return x.Key
}

func (x *MessageV1) Check() error {
	vs := x.generateValidSum()
	if x.ValidSum != vs {
		return ErrMessageIsInvalid
	}
	return nil
}

func (x *MessageV1) Marshal() ([]byte, error) {
	return proto.Marshal((*messagev1.Message)(x))
}

func (x *MessageV1) UnMarshal(bytes []byte) error {
	if x == nil {
		return ErrMessageIsNil
	}
	return proto.Unmarshal(bytes, (*messagev1.Message)(x))
}

// ------------------------------ private ------------------------------

func (x *MessageV1) generateUniKey() string {
	b, _ := x.Marshal()
	c32ID := crc32.ChecksumIEEE(b)
	return fmt.Sprintf("%d-%d", time.Now().Unix(), c32ID)
}

func (x *MessageV1) generateValidSum() uint32 {
	m := *x
	m.ValidSum = 0
	b, _ := m.Marshal()
	return crc32.ChecksumIEEE(b)
}

func (x *MessageV1) assignValidSum() {
	x.ValidSum = 0
	x.ValidSum = x.generateValidSum()
}
