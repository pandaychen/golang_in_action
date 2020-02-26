package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"golang.org/x/net/context"
	"io"
	"time"
)

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.0.0.1:6831",
		},
	}
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func foo3(req string, ctx context.Context) (reply string) {
	//1.创建子span
	span, _ := opentracing.StartSpanFromContext(ctx, "span_foo3")
	defer func() {
		//4.接口调用完，在tag中设置request和reply
		span.SetTag("request", req)
		span.SetTag("reply", reply)
		span.Finish()
	}()

	println(req)
	//2.模拟处理耗时
	time.Sleep(time.Second / 2)
	//3.返回reply
	reply = "foo3Reply"
	return
}

//跟foo3一样逻辑
func foo4(req string, ctx context.Context) (reply string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "span_foo4")
	defer func() {
		span.SetTag("request", req)
		span.SetTag("reply", reply)
		span.Finish()
	}()

	println(req)
	time.Sleep(time.Second / 2)
	reply = "foo4Reply"
	return
}

func main() {
	tracer, closer := initJaeger("jaeger-demo")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer) //StartspanFromContext创建新span时会用到

	span := tracer.StartSpan("span_root")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	r1 := foo3("Hello foo3", ctx)
	r2 := foo4("Hello foo4", ctx)
	fmt.Println(r1, r2)
	span.Finish()
}
