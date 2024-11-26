# pulsartracing
封装pulsar库以支持 opentelemetry.

[![Go](https://github.com/me-cs/pulsartracing/workflows/Go/badge.svg)](https://github.com/me-cs/pulsartracing/actions)
[![codecov](https://codecov.io/gh/me-cs/pulsartracing/branch/main/graph/badge.svg)](https://codecov.io/gh/me-cs/pulsartracing)
[![Release](https://img.shields.io/github/v/release/me-cs/pulsartracing.svg?style=flat-square)](https://github.com/me-cs/pulsartracing)
[![Go Report Card](https://goreportcard.com/badge/github.com/me-cs/pulsartracing)](https://goreportcard.com/report/github.com/me-cs/pulsartracing)
[![Go Reference](https://pkg.go.dev/badge/github.com/me-cs/pulsartracing.svg)](https://pkg.go.dev/github.com/me-cs/pulsartracing)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

简体中文 | [English](README.md)

### 使用示例:
```go
package main

import (
	"context"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/me-cs/pulsartracing"
)

var pulsarClient pulsar.Client

func produce() {
	p, err := pulsarClient.CreateProducer(pulsar.ProducerOptions{})
	if err != nil {
		panic(err)
	}
	//一般来说这里需要传递已经被跟踪过的上游context
	_, _ = p.Send(context.Background(), &pulsar.ProducerMessage{})
}

func consume() {
	customerConsumer, err := pulsarClient.Subscribe(pulsar.ConsumerOptions{})
	if err != nil {
		panic(err)
	}
	for {
		ctx := context.Background()
		ctx, msg, err := pulsartracing.ReceiveWithSpanCtx(ctx, customerConsumer)
		if err != nil {
			continue
		}
		err = customerConsumer.Ack(msg)
		if err != nil {
			continue
		}
		//把这个ctx继续传递给下游
		//downstream(ctx)
		//然后你就可以在你的链路追踪系统（例如jaeger）里看到消息被跟踪到了
	}
}

func main() {
	var err error
	pulsarClient, err = pulsartracing.NewClient(pulsar.ClientOptions{
		URL: "pulsar://pulsar.xxx.cn:6650",
	})
	if err != nil {
		panic(err)
	}
	produce()
	consume()
}

```

更重要的
假设你的应用程序是这样初始化opentelemetry的，你需要把opentelemetry桥接到opentracing
```go
package main

import (
	"github.com/opentracing/opentracing-go"
	"go.opentelemetry.io/otel"
	otelBridge "go.opentelemetry.io/otel/bridge/opentracing"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	tp *sdktrace.TracerProvider
)

func main() {
	opts := []sdktrace.TracerProviderOption{}
	tp = sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otelTracer := tp.Tracer("you trace name")
	// Use the bridgeTracer as your OpenTracing tracer.
	bridgeTracer, wrapperTracerProvider := otelBridge.NewTracerPair(otelTracer)
	// Set the wrapperTracerProvider as the global OpenTelemetry
	// TracerProvider so instrumentation will use it by default.
	otel.SetTracerProvider(wrapperTracerProvider)
	opentracing.SetGlobalTracer(bridgeTracer)
	return
}

```