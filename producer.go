package pulsartracing

import (
	"context"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/me-cs/pulsartracing/tracing"
)

type producer struct {
	pulsar.Producer
}

func (p *producer) Send(ctx context.Context, message *pulsar.ProducerMessage) (pulsar.MessageID, error) {
	tracing.InjectProducerMessageSpanContext(ctx, message)
	return p.Producer.Send(ctx, message)
}

func (p *producer) SendAsync(ctx context.Context, message *pulsar.ProducerMessage, cb func(pulsar.MessageID, *pulsar.ProducerMessage, error)) {
	tracing.InjectProducerMessageSpanContext(ctx, message)
	p.Producer.SendAsync(ctx, message, cb)
}
