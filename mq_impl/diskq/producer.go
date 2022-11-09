package diskq

import (
	"context"
	"sync"
	"time"

	"github.com/Ccheers/kratos-mq/mq"
	"github.com/Ccheers/kratos-mq/mq_impl/diskq/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nsqio/go-diskqueue"
	"golang.org/x/sync/singleflight"
)

type ProducerImpl struct {
	sf       singleflight.Group
	mu       sync.Mutex
	producer map[string]diskqueue.Interface
	c        *config.Config
	logger   log.Logger
}

func NewProducer(c *config.Config, logger log.Logger) (mq.Producer, error) {
	// Instantiate a producer.
	return &ProducerImpl{
		producer: make(map[string]diskqueue.Interface),
		c:        c,
		logger:   logger,
	}, nil
}

func (x *ProducerImpl) Publish(ctx context.Context, topic string, message mq.Message) error {
	b, err := message.Marshal()
	if err != nil {
		return err
	}
	ins, _, _ := x.sf.Do(topic, func() (interface{}, error) {
		x.mu.Lock()
		defer x.mu.Unlock()
		if ins, ok := x.producer[topic]; ok {
			return ins, nil
		}
		queue := gDiskQueueManager.NewDiskQueue(topic, x.c.DataPath, x.c.MaxBytesPerFile, x.c.MinMsgSize, x.c.MaxMsgSize, x.c.SyncEvery, time.Duration(x.c.SyncTimeout)*time.Millisecond, x.logger)
		x.producer[topic] = queue
		return queue, nil
	})
	return ins.(diskqueue.Interface).Put(b)
}

func (x *ProducerImpl) Close(ctx context.Context) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	var err error
	for _, producer := range x.producer {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = producer.Close()
		if err != nil {
			x.logger.Log(log.LevelError, "module", "[diskq][Producer][Close]", "err", err)
		}
	}
	return nil
}
