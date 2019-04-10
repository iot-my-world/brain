package sync

import (
	"gitlab.com/iotTracker/brain/log"
	messagingProducer "gitlab.com/iotTracker/brain/messaging/producer"
	producerException "gitlab.com/iotTracker/brain/messaging/producer/exception"
	"gopkg.in/Shopify/sarama.v1"
)

type producer struct {
	producer sarama.SyncProducer
	brokers  []string
	topic    string
}

func New(
	brokers []string,
	topic string,
) messagingProducer.Producer {
	return &producer{
		brokers: brokers,
		topic:   topic,
	}
}

func (p *producer) Start() error {
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewSyncProducer(p.brokers, config)
	if err != nil {
		return producerException.Start{Reasons: []string{"failed to connect new producer", err.Error()}}
	}

	p.producer = producer

	return nil
}

func (p *producer) Produce(data []byte) error {
	// We are not setting a message key, which means that all messages will
	// be distributed randomly over the different partitions.

	//partition, offset, err := p.producer.SendMessage(&sarama.ProducerMessage{
	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(data),
	})
	if err != nil {
		return producerException.Produce{Reasons: []string{err.Error()}}
	} else {
		// The tuple (topic, partition, offset) can be used as a unique identifier
		// for a message in a Kafka cluster.
		log.Debug("Published kafka message")
	}
	return nil
}
