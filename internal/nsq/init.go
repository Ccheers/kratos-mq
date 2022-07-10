package nsq

import (
	"time"

	"github.com/Ccheers/kratos-mq/internal/nsq/config"
	"github.com/google/wire"
	"github.com/nsqio/go-nsq"
)

var ProviderSet = wire.NewSet(NewConsumer, NewProducer)

func NewNsqConfigWithConfig(c *config.Config) *nsq.Config {
	// Instantiate a consumer that will subscribe to the provided channel.
	cfg := nsq.NewConfig()
	cfg.DialTimeout = time.Second * time.Duration(c.DialTimeout)
	cfg.ReadTimeout = time.Second * time.Duration(c.ReadTimeout)
	cfg.WriteTimeout = time.Second * time.Duration(c.WriteTimeout)
	return cfg
}
