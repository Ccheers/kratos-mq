package diskq

import (
	"context"
	"fmt"
	"sync"
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
	consumerMap  map[string][]chan mq.Message
	consumerChan map[string]chan mq.Message
	pool         routinepool.Pool
}

func NewConsumer(c *config.Config, logger log.Logger) (mq.Consumer, error) {
	// Instantiate a consumer that will subscribe to the provided channel.
	return &ConsumerImpl{
		c:            c,
		logger:       logger,
		consumerChan: make(map[string]chan mq.Message),
		consumerMap:  make(map[string][]chan mq.Message),
		pool:         routinepool.NewPool("[diskq][Consumer]", 4, routinepool.NewConfig()),
	}, nil
}

func (x *ConsumerImpl) Subscribe(ctx context.Context, topic string, channel string) (<-chan mq.Message, error) {
	x.mu.Lock()
	defer x.mu.Unlock()

	uniKey := x.generateKey(topic, channel)
	if ch, ok := x.consumerChan[uniKey]; ok {
		return ch, nil
	}

	_, _, _ = x.sf.Do(topic, func() (interface{}, error) {
		if consumer, ok := x.consumerMap[topic]; ok {
			return consumer, nil
		}
		queue := gDiskQueueManager.NewDiskQueue(topic, x.c.DataPath, x.c.MaxBytesPerFile, x.c.MinMsgSize, x.c.MaxMsgSize, x.c.SyncEvery, time.Duration(x.c.SyncTimeout)*time.Millisecond, x.logger)
		x.pool.Go(func(ctx context.Context) {
			for {
				select {
				case body := <-queue.ReadChan():
					msg, err := mq.NewMessageFromByte(body)
					if err != nil {
						x.logger.Log(log.LevelError, "module", "NewMessageFromByte", "err", err, "body", string(body))
					}
					for _, cch := range x.consumerMap[topic] {
						select {
						case cch <- msg:
						default:
						}
					}
				}
			}
		})
		return queue, nil
	})

	ch := make(chan mq.Message, 1024)
	x.consumerMap[topic] = append(x.consumerMap[topic], ch)
	x.consumerChan[uniKey] = ch
	return ch, nil
}

func (x *ConsumerImpl) Close(ctx context.Context) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	for uniKey, _ := range x.consumerChan {
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
