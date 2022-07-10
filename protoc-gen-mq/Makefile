gen_example:
	go install
	protoc --proto_path=. \
           --proto_path=./example/api \
           --go_out=paths=source_relative:. \
           --go-grpc_out=paths=source_relative:. \
           --mq_out=paths=source_relative:. \
           example/api/product/app/v1/v1.proto
	#protoc-go-inject-tag -input=./example/api/product/app/v1/v1.pb.go