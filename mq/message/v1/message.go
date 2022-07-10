package v1

import (
	"errors"
	"fmt"
	"hash/crc32"
	"time"

	"github.com/Ccheers/kratos-mq/mq"
	"github.com/go-kratos/kratos/v2/metadata"
	"google.golang.org/protobuf/proto"
)

var (
	ErrMessageIsNil     = errors.New("message is nil show alloc memory")
	ErrMessageIsInvalid = errors.New("message is invalid")
)

type MessageOption func(x *Message)

func WithMetadata(md metadata.Metadata) MessageOption {
	return func(x *Message) {
		x.Md = md
	}
}

var _ mq.Message = (*Message)(nil)

func NewMessage(payload mq.Payload, opts ...MessageOption) mq.Message {
	m := &Message{Data: payload}
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

func NewMessageFromByte(b []byte) (mq.Message, error) {
	m := &Message{}
	err := m.UnMarshal(b)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (x *Message) Metadata() metadata.Metadata {
	return x.Md
}

func (x *Message) Payload() mq.Payload {
	return x.Data
}

func (x *Message) Err() error {
	return fmt.Errorf(x.Error)
}

func (x *Message) UniKey() string {
	return x.Key
}

func (x *Message) Check() error {
	vs := x.generateValidSum()
	if x.ValidSum != vs {
		return ErrMessageIsInvalid
	}
	return nil
}

func (x *Message) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Message) UnMarshal(bytes []byte) error {
	if x == nil {
		return ErrMessageIsNil
	}
	return proto.Unmarshal(bytes, x)
}

// ------------------------------ private ------------------------------

func (x *Message) generateUniKey() string {
	b, _ := x.Marshal()
	c32ID := crc32.ChecksumIEEE(b)
	return fmt.Sprintf("%d-%d", time.Now().Unix(), c32ID)
}

func (x *Message) generateValidSum() uint32 {
	m := *x
	m.ValidSum = 0
	b, _ := m.Marshal()
	return crc32.ChecksumIEEE(b)
}

func (x *Message) assignValidSum() {
	x.ValidSum = 0
	x.ValidSum = x.generateValidSum()
}
