package nsq

import (
	"time"

	"github.com/Ccheers/kratos-mq/mq_impl/nsq/config"
	"github.com/nsqio/go-nsq"
)

func NewNsqConfigWithConfig(c *config.Config) *nsq.Config {
	// Instantiate a consumer that will subscribe to the provided channel.
	cfg := nsq.NewConfig()
	cfg.DialTimeout = time.Second * time.Duration(c.DialTimeout)
	cfg.ReadTimeout = time.Second * time.Duration(c.ReadTimeout)
	cfg.WriteTimeout = time.Second * time.Duration(c.WriteTimeout)
	cfg.HeartbeatInterval = time.Second * time.Duration(c.HeartbeatInterval)
	cfg.MaxInFlight = int(c.BatchSize)
	return cfg
}
