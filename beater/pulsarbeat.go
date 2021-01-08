package beater

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/yukshimizu/pulsarbeat/config"
	"sync"
	"time"
)

// pulsarbeat configuration.
type pulsarbeat struct {
	done            chan struct{}
	config          config.Config
	client          beat.Client
	pulsarClient    *pulsar.Client
	pulsarConsumers *[]pulsar.Consumer
}

const selector string = "pulsarbeat"

// New creates an instance of pulsarbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	logp.Debug(selector, "Default config yml is: %#v", c)
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}
	logp.Debug(selector, "After reading config yml is: %#v", c)

	client, err := config.NewPulsarClient(c.Client)
	if err != nil {
		return nil, fmt.Errorf("error creating pulsar client: %v", err)
	}

	consumers, err := config.NewPulsarConsumer(client, c.Consumer)
	if err != nil {
		return nil, fmt.Errorf("error creating pulsar consumer: %v", err)
	}

	bt := &pulsarbeat{
		done:            make(chan struct{}),
		config:          c,
		pulsarClient:    client,
		pulsarConsumers: consumers,
	}

	return bt, nil
}

// Run starts pulsarbeat.
func (bt *pulsarbeat) Run(b *beat.Beat) error {
	logp.Info("pulsarbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		(*bt.pulsarClient).Close()
		logp.Debug(selector, "pulsar client Closed!")
	}()

	for _, consumer := range *bt.pulsarConsumers {
		wg.Add(1)
		consumer := consumer
		go func(ctx context.Context) {
			defer wg.Done()
			defer func() {
				consumer.Close()
				logp.Debug(selector, "pulsar consumer: %#v Closed!", consumer)
			}()

			for {
				select {
				case <-ctx.Done():
					logp.Debug(selector, "done ctx")
					return
				default:
					// pulsar normal implementation
					msg, err := consumer.Receive(ctx)
					if err != nil {
						logp.Debug(selector, "consumer Receive failed: %v", err)
						if err != context.Canceled {
							consumer.Nack(msg)
						}
					} else {

						logp.Debug(selector, "Received message msgId: %#v -- content: '%s'",
							msg.ID(), string(msg.Payload()))

						event := beat.Event{
							Timestamp: time.Now(),
							Fields: common.MapStr{
								"pulsar": common.MapStr{
									"topic":     msg.Topic(),
									"producer":  msg.ProducerName(),
									"key":       msg.Key(),
									"timestamp": msg.PublishTime(),
								},
								"message": string(msg.Payload()),
							},
						}

						bt.client.Publish(event)
						consumer.Ack(msg)

						logp.Debug(selector, "Event sent")
					}
				}
			}
		}(ctx)
	}

	go func() {
		<-bt.done
		logp.Debug(selector, "calling cancel function")
		cancel()
	}()

	wg.Wait()
	logp.Debug(selector, "Run method exited!")
	return nil
}

// Stop stops pulsarbeat.
func (bt *pulsarbeat) Stop() {
	logp.Debug(selector, "Stop method called")
	bt.client.Close()
	close(bt.done)
}
