module github.com/Ccheers/kratos-mq/protoc-gen-gin/example

go 1.16

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/go-kratos/kratos/v2 v2.2.1
	google.golang.org/genproto v0.0.0-20220126215142-9970aeb2e350
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
	github.com/Ccheers/kratos-mq v0.0.1
)

replace (
	github.com/Ccheers/kratos-mq => ../../
)