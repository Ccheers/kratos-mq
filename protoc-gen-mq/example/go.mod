module github.com/Ccheers/kratos-mq/protoc-gen-mq/example

go 1.16

require (
	github.com/Ccheers/kratos-mq v0.0.1
	github.com/go-kratos/kratos/v2 v2.3.1
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.0
)

replace github.com/Ccheers/kratos-mq => ../../
