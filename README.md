# pulsartracing
Wrap pulsar so that it supports opentelemetry.

[![Go](https://github.com/me-cs/pulsartracing/workflows/Go/badge.svg)](https://github.com/me-cs/pulsartracing/actions)
[![codecov](https://codecov.io/gh/me-cs/pulsartracing/branch/main/graph/badge.svg)](https://codecov.io/gh/me-cs/pulsartracing)
[![Release](https://img.shields.io/github/v/release/me-cs/pulsartracing.svg?style=flat-square)](https://github.com/me-cs/pulsartracing)
[![Go Report Card](https://goreportcard.com/badge/github.com/me-cs/pulsartracing)](https://goreportcard.com/report/github.com/me-cs/pulsartracing)
[![Go Reference](https://pkg.go.dev/badge/github.com/me-cs/pulsartracing.svg)](https://pkg.go.dev/github.com/me-cs/pulsartracing)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

English | [简体中文](README-CN.md)

### Example use:

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
	//Generally an upstream context that has already been traced.
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
		//Pass this context to the downstream
		//downstream(ctx)
		//Then you can see in your link tracking system (e.g. jaeger) that the message was tracked to
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


More Important to Note
Assuming your application initializes opentelemetry like this, 
you need to bridge opentelemetry to opentracing
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