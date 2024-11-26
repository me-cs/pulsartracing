package pulsarMq

import (
	"context"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/me-cs/pulsartracing/tracing"
	"github.com/opentracing/opentracing-go"
	otelBridge "go.opentelemetry.io/otel/bridge/opentracing"
)

type consumer struct {
	pulsar.Consumer
}

// Receive
// use ReceiveWithSpanCtx if you need tracing message
func (c *consumer) Receive(ctx context.Context) (pulsar.Message, error) {
	return c.Consumer.Receive(ctx)
}

func ReceiveWithSpanCtx(ctx context.Context, c pulsar.Consumer) (context.Context, pulsar.Message, error) {
	msg, err := c.Receive(ctx)
	if err != nil {
		return nil, nil, err
	}
	cm := pulsar.ConsumerMessage{
		Consumer: c,
		Message:  msg,
	}
	span := tracing.CreateSpanFromMessage(&cm, opentracing.GlobalTracer(), "child_span")
	bridge := otelBridge.NewBridgeTracer()
	ctx = bridge.ContextWithSpanHook(ctx, span)
	span.Finish()
	return ctx, msg, nil
}
