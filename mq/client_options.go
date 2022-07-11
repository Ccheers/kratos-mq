package mq

type ClientOptionFunc func(*clientOptions)

// ClientOptionWithEncodeFunc 中间件
func ClientOptionWithEncodeFunc(f EncodeFunc) ClientOptionFunc {
	return func(options *clientOptions) {
		options.encodeFunc = f
	}
}
