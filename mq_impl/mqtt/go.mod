module github.com/Ccheers/kratos-mq/mq_impl/mqtt

go 1.17

require (
	github.com/Ccheers/kratos-mq v0.0.0-00010101000000-000000000000
	github.com/eclipse/paho.mqtt.golang v1.4.1
	github.com/go-kratos/kratos/v2 v2.3.1
)

require (
	github.com/Ccheers/kratos-mq/protoc-gen-mq v0.0.0-00010101000000-000000000000 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/ccheers/xpkg v1.0.2 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sync v0.0.0-20220513210516-0976fa681c29 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd // indirect
	google.golang.org/grpc v1.46.2 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/Ccheers/kratos-mq => ../../

replace github.com/Ccheers/kratos-mq/protoc-gen-mq => ../../protoc-gen-mq
