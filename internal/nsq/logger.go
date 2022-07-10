package nsq

import "github.com/go-kratos/kratos/v2/log"

type logWrapper struct {
	l log.Logger
}

func newLogWrapper(l log.Logger) *logWrapper {
	return &logWrapper{l: l}
}

func (x *logWrapper) Output(calldepth int, s string) error {
	return x.l.Log(log.LevelInfo, "msg", s)
}
