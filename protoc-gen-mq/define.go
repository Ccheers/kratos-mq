package main

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

type service struct {
	Name      string // Greeter
	FullName  string // helloworld.Greeter
	FilePath  string // api/helloworld/helloworld.proto
	Comment   string // 注释
	Methods   []*method
	MethodSet map[string]*method
}

// InterfaceName service interface name
func (s *service) InterfaceName() string {
	return s.Name + "MQServer"
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

// matchMQReg 匹配 mq tag @mq: topic::channel
var matchMQReg = regexp.MustCompile("@mq:\\s*([\\w\\\\.]+)::([\\w\\\\.]+)")

func genMethod(m *protogen.Method, g *protogen.GeneratedFile) []*method {
	var methods []*method

	mqs := matchMQReg.FindAllString(m.Comments.Leading.String(), -1)
	for _, str := range mqs {
		str = str[4:]
		s := strings.Split(str, "::")
		topic := strings.TrimSpace(s[0])
		channel := strings.TrimSpace(s[1])
		methods = append(methods, buildMethodDesc(m, g, topic, channel))
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
