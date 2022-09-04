package v1

//go:generate protoc --proto_path=. --proto_path ../../third_party --go_out=paths=source_relative:. --mq_out=paths=source_relative:. ./hello.proto
