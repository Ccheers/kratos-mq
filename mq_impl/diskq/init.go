package diskq

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/nsqio/go-diskqueue"
	"golang.org/x/sync/singleflight"
)

var gDiskQueueManager = &diskQueueManager{
	queueMap: make(map[string]diskqueue.Interface),
	logger:   log.DefaultLogger,
}

type diskQueueManager struct {
	queueMap map[string]diskqueue.Interface
	sf       singleflight.Group
	mu       sync.Mutex
	close    uint32
	logger   log.Logger
}

func (x *diskQueueManager) Close(ctx context.Context) error {
	if atomic.CompareAndSwapUint32(&x.close, 0, 1) {
		return nil
	}
	for _, queue := range x.queueMap {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		err := queue.Close()
		x.logger.Log(log.LevelError, "module", "[diskQueueManager][Close]", "err", err)
	}
	return nil
}

func (x *diskQueueManager) SetLogger(logger log.Logger) {
	x.logger = logger
}

func (x *diskQueueManager) NewDiskQueue(topic string, dataPath string, maxBytesPerFile int64, minMsgSize int32, maxMsgSize int32, syncEvery int64, syncTimeout time.Duration, logger log.Logger) diskqueue.Interface {
	res, _, _ := x.sf.Do(topic, func() (interface{}, error) {
		x.mu.Lock()
		defer x.mu.Unlock()
		if ins, ok := x.queueMap[topic]; ok {
			return ins, nil
		}
		queue := diskqueue.New(
			topic,
			dataPath,
			maxBytesPerFile,
			minMsgSize,
			maxMsgSize,
			syncEvery,
			syncTimeout,
			func(lvl diskqueue.LogLevel, f string, args ...interface{}) {
				switch lvl {
				case diskqueue.DEBUG:
					logger.Log(log.LevelDebug, args...)
				case diskqueue.INFO:
					logger.Log(log.LevelInfo, args...)
				case diskqueue.WARN:
					logger.Log(log.LevelWarn, args...)
				case diskqueue.ERROR:
					logger.Log(log.LevelError, args...)
				case diskqueue.FATAL:
					logger.Log(log.LevelFatal, args...)
				}
			},
		)
		x.queueMap[topic] = queue
		return queue, nil
	})
	return res.(diskqueue.Interface)
}
