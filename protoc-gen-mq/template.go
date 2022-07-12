package main

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	contextPkg        = protogen.GoImportPath("context")
	fmtPkg            = protogen.GoImportPath("fmt")
	kMQPkg            = protogen.GoImportPath("github.com/Ccheers/kratos-mq/mq")
	kratosEncodingPkg = protogen.GoImportPath("github.com/go-kratos/kratos/v2/encoding")
	middlewarePkg     = protogen.GoImportPath("github.com/go-kratos/kratos/v2/middleware")
	metadataPkg       = protogen.GoImportPath("github.com/go-kratos/kratos/v2/metadata")

	deprecationComment = "// Deprecated: Do not use."
)

var methodSets = make(map[string]int)

// generateFile generates a _gin.pb.go file.
func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + "_kmq.kratos.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by github.com/Ccheers/kratos-mq/protoc-gen-mq. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.P("// This is a compile-time assertion to ensure that this generated file")
	g.P("// is compatible with the Ccheers/protoc-gen-mq package it is being compiled against.")
	g.P()

	for _, service := range file.Services {
		genService(gen, file, g, service)
	}
	return g
}

func genService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, s *protogen.Service) {
	if s.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}

	// HTTP Server.
	sd := &service{
		Name:     s.GoName,
		FullName: string(s.Desc.FullName()),
		FilePath: file.Desc.Path(),
		Comment:  clearComment(string(s.Comments.Leading)),
	}

	for _, method := range s.Methods {
		sd.Methods = append(sd.Methods, genMethod(method, g)...)
	}
	if sd.MethodSet == nil {
		sd.MethodSet = map[string]*method{}
		for _, m := range sd.Methods {
			m := m
			sd.MethodSet[m.Name] = m
		}
	}
	// 生成接口
	generateServerInterface(g, sd)
	// 生成注册函数
	generateServerRegisterFunc(g, sd)
	// 生成函数列表
	generateServerMethodList(g, sd)
	// 生成接口
	generateClientInterface(g, sd)
	// 生成函数列表
	generateClientMethodList(g, sd)
}

// 生成接口
func generateServerInterface(g *protogen.GeneratedFile, s *service) {
	g.P("// ", s.Comment)
	g.P("type ", s.ServerInterfaceName(), " interface {")
	for _, m := range s.MethodSet {
		g.P("\t", generateServerMethodName(m.Name), "(", contextPkg.Ident("Context"), ", *", m.Request, ") error")
	}
	g.P("}")
}

// 生成注册函数
func generateServerRegisterFunc(g *protogen.GeneratedFile, s *service) {
	g.P("func Register", s.ServerInterfaceName(), "(svr *", kMQPkg.Ident("Server"), ", srv ", s.ServerInterfaceName(), ") error {")
	g.P("var err error")
	for _, m := range s.Methods {
		g.P("\t", "err = svr.Subscriber(\"", m.Topic, "\", \"", m.Channel, "\", ", m.HandlerName(), "(svr, srv))")
		{
			g.P("if err != nil {")
			g.P("return err")
			g.P("}")
		}
	}
	g.P("return nil")
	g.P("}")
}

func generateServerMethodName(name string) string {
	return fmt.Sprintf("MQ_%s", name)
}

// generateServerMethodList 函数列表
func generateServerMethodList(g *protogen.GeneratedFile, s *service) {
	for _, m := range s.Methods {
		g.P("func ", m.HandlerName(), "(svr *", kMQPkg.Ident("Server"), ", srv ", s.ServerInterfaceName(), ") ", kMQPkg.Ident("HandleFunc"), " {")
		g.P("return func(ctx ", contextPkg.Ident("Context"), ", message ", kMQPkg.Ident("Message"), ") {")
		g.P("")

		// 声明参数
		{
			g.P("var in ", m.Request)
			g.P("var err error")
			g.P("")
		}
		// 编码解码器
		{
			g.P("err = svr.DecodeFunc()(ctx, message, &in)")
			g.P("if err != nil {")
			g.P("svr.ErrHandler().Handle(err)")
			g.P("return")
			g.P("}")
			g.P("")
		}
		// 自制 metadata
		{
			g.P("md := ", metadataPkg.Ident("New"), "(nil)")
			g.P("")

			g.P("for k, v := range message.Metadata() {")
			g.P("md.Set(k, v)")
			g.P("}")
			g.P("")
			g.P("newCtx := ", metadataPkg.Ident("NewServerContext"), "(ctx, md)")
			g.P("")
		}

		// 加入中间件能力
		{
			g.P("ms := ", kMQPkg.Ident("GetMiddlewareFromContext(ctx)"))
			g.P("")

			g.P("handler := func(ctx ", contextPkg.Ident("Context"), ", req interface{}) (interface{}, error) {")
			g.P("err := srv.", generateServerMethodName(m.Name), "(ctx, req.(*", m.Request, "))")
			g.P("return nil, err")
			g.P("}")
			g.P("")
		}
		// 执行请求
		{
			g.P("_, err = ", middlewarePkg.Ident("Chain"), "(ms...)(handler)(newCtx, &in)")
			g.P("if err != nil {")
			g.P("svr.ErrHandler().Handle(err)")
			g.P("return")
			g.P("}")
			g.P("")
		}
		g.P("}")
		g.P("}")
	}
}

// 生成接口
func generateClientInterface(g *protogen.GeneratedFile, s *service) {
	g.P("type ", s.ClientInterfaceName(), " interface {")
	for _, m := range s.Methods {
		g.P("\t", fmt.Sprintf("%s_%s", m.Name, m.Topic), "(", contextPkg.Ident("Context"), ", *", m.Request, ") error")
	}
	g.P("}")
}

// generateClientMethodList 函数列表
func generateClientMethodList(g *protogen.GeneratedFile, s *service) {
	// 客户端接口定义
	clientStructName := fmt.Sprintf("%sImpl", s.ClientInterfaceName())
	{
		g.P("type ", clientStructName, " struct {")
		g.P("\t", "cc *", kMQPkg.Ident("Client"))
		g.P("}")
	}
	// 构造函数
	g.P("func ", fmt.Sprintf("New%s", s.ClientInterfaceName()), "(cc *", kMQPkg.Ident("Client"), ") ", s.ClientInterfaceName(), " {")
	g.P("\t", "return &", clientStructName, "{cc: cc}")
	g.P("}")

	// 函数列表
	for _, m := range s.Methods {
		g.P(fmt.Sprintf("func (x *%s)", clientStructName), fmt.Sprintf("%s_%s", m.Name, m.Topic), "(ctx ", contextPkg.Ident("Context"), ", req *", m.Request, ") error {")
		g.P("return x.cc.Invoke(ctx, \"", m.Topic, "\", req)")
		g.P("}")
	}
}
