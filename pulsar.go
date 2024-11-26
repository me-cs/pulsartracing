package pulsarMq

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/me-cs/pulsartracing/tracing"
)

type client struct {
	pulsar.Client
}

func NewClient(clientOpt pulsar.ClientOptions) (pulsar.Client, error) {
	c, err := pulsar.NewClient(clientOpt)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

func (c *client) CreateProducer(opt pulsar.ProducerOptions) (pulsar.Producer, error) {
	opt.Interceptors = append(opt.Interceptors, &tracing.ProducerInterceptor{})
	p, err := c.Client.CreateProducer(opt)
	return &producer{p}, err
}

func (c *client) Subscribe(opt pulsar.ConsumerOptions) (pulsar.Consumer, error) {
	opt.Interceptors = append(opt.Interceptors, &tracing.ConsumerInterceptor{})
	cc, err := c.Client.Subscribe(opt)
	return &consumer{cc}, err
}
