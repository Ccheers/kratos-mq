package main

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/Ccheers/kratos-mq/protoc-gen-mq/proto/kmq/v1/options"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

type service struct {
	Name      string // Greeter
	FullName  string // helloworld.Greeter
	FilePath  string // api/helloworld/helloworld.proto
	Comment   string // 注释
	Methods   []*method
	MethodSet map[string]*method
}

// ServerInterfaceName server interface name
func (s *service) ServerInterfaceName() string {
	return s.Name + "MQServer"
}

// ClientInterfaceName client interface name
func (s *service) ClientInterfaceName() string {
	return s.Name + "MQClient"
}

type method struct {
	Name    string // SayHello
	Num     int    // 一个 rpc 方法可以对应多个 http 请求
	Request string // SayHelloReq
	Reply   string // SayHelloResp
	Comment string // 注释
	// mq
	Topic   string // topic
	Channel string // channel
}

// HandlerName for gin handler name
func (m *method) HandlerName() string {
	return fmt.Sprintf("%s_%d_MQHandler", m.Name, m.Num)
}

func genMethod(m *protogen.Method, g *protogen.GeneratedFile) []*method {
	var methods []*method

	mq, ok := proto.GetExtension(m.Desc.Options(), options.E_Mq).(*options.MQ)
	if !ok {
		return nil
	}
	for _, group := range mq.GetSubscribes() {
		methods = append(methods, buildMethodDesc(m, g, group.GetTopic(), group.GetChannel()))
	}

	return methods
}

func buildMethodDesc(m *protogen.Method, g *protogen.GeneratedFile, topic, channel string) *method {
	defer func() { methodSets[m.GoName]++ }()
	md := &method{
		Name:    m.GoName,
		Num:     methodSets[m.GoName],
		Request: g.QualifiedGoIdent(m.Input.GoIdent),
		Reply:   g.QualifiedGoIdent(m.Output.GoIdent),
		Topic:   topic,
		Channel: channel,
		Comment: clearComment(string(m.Comments.Leading)),
	}
	return md
}

func clearComment(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, "\n", ""))
}
