package mq

import (
	"context"

	"github.com/go-kratos/kratos/v2/metadata"
)

type Payload []byte

// Message 消息定义
type Message interface {
	// Metadata 元数据
	Metadata() metadata.Metadata
	// Payload 荷载信息
	Payload() Payload
	// Err 错误
	Err() error
	// UniKey 消息唯一键 用于消息去重(已经消费的消息不再消费)
	UniKey() string

	// Check 检查消息完整性
	Check() error
	// Marshal 序列化消息
	Marshal() ([]byte, error)
	// UnMarshal 解析消息
	UnMarshal([]byte) error
}

type Consumer interface {
	// Subscribe topic 是主题，相同主题不同 channel 可以分发相同消息实现广播
	Subscribe(ctx context.Context, topic string, channel string) (<-chan Message, error)
	// Close 停止消费
	Close(ctx context.Context) error
}

type Producer interface {
	// Publish topic 是主题 向 topic 投递消息，监听此 topic 的 Producer 可以收到消息
	Publish(ctx context.Context, topic string, message Message) error
	// Close 停止消费
	Close(ctx context.Context) error
}
