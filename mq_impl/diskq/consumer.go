package diskq

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Ccheers/kratos-mq/mq"
	"github.com/Ccheers/kratos-mq/mq_impl/diskq/config"
	"github.com/ccheers/xpkg/sync/routinepool"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/singleflight"
)

var _ mq.Consumer = (*ConsumerImpl)(nil)

type ConsumerImpl struct {
	c *config.Config

	logger log.Logger

	mu           sync.Mutex
	sf           singleflight.Group
	consumerChan map[string]chan mq.Message
	pool         routinepool.Pool

	status uint32
}

func NewConsumer(c *config.Config, logger log.Logger) (mq.Consumer, error) {
	// Instantiate a consumer that will subscribe to the provided channel.
	return &ConsumerImpl{
		c:            c,
		logger:       logger,
		consumerChan: make(map[string]chan mq.Message),
		pool:         routinepool.NewPool("[diskq][Consumer]", 4, routinepool.NewConfig()),
		status:       statusRunning,
	}, nil
}

func (x *ConsumerImpl) Subscribe(ctx context.Context, topic string, channel string) (<-chan mq.Message, error) {
	x.mu.Lock()
	defer x.mu.Unlock()

	if ch, ok := x.consumerChan[topic]; ok {
		return ch, nil
	}

	ch, _, _ := x.sf.Do(topic, func() (interface{}, error) {
		if ch, ok := x.consumerChan[topic]; ok {
			return ch, nil
		}
		queue := gDiskQueueManager.NewDiskQueue(topic, x.c.DataPath, x.c.MaxBytesPerFile, x.c.MinMsgSize, x.c.MaxMsgSize, x.c.SyncEvery, time.Duration(x.c.SyncTimeout)*time.Millisecond, x.logger)
		ch := make(chan mq.Message, 1)
		x.pool.Go(func(ctx context.Context) {
			for {
				if atomic.LoadUint32(&x.status) == statusClosed {
					return
				}
				select {
				case body := <-queue.ReadChan():
					msg, err := mq.NewMessageFromByte(body)
					if err != nil {
						_ = x.logger.Log(log.LevelError, "module", "NewMessageFromByte", "err", err, "body", string(body))
					}
					ch <- msg
				}
			}
		})
		return ch, nil
	})

	x.consumerChan[topic] = ch.(chan mq.Message)
	return x.consumerChan[topic], nil
}

func (x *ConsumerImpl) Close(ctx context.Context) error {
	if !atomic.CompareAndSwapUint32(&x.status, statusRunning, statusClosed) {
		return nil
	}
	x.mu.Lock()
	defer x.mu.Unlock()
	for uniKey := range x.consumerChan {
		close(x.consumerChan[uniKey])
	}
	err := gDiskQueueManager.Close(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (x *ConsumerImpl) generateKey(topic string, channel string) string {
	return fmt.Sprintf("%s::%s", topic, channel)
}
