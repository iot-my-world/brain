package group

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"gitlab.com/iotTracker/brain/log"
	consumerGroupException "gitlab.com/iotTracker/brain/messaging/consumer/group/exception"
	"os"
	"os/signal"
	"syscall"
)

// consumer represents a Sarama consumer group consumer
type consumer struct {
	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	c.ready <- true
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		log.Info(fmt.Sprintf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic))
		session.MarkMessage(message, "")
	}

	return nil
}

type group struct {
	brokers   []string
	topics    []string
	groupName string
}

func New(
	brokers []string,
	topics []string,
	groupName string,
) *group {
	return &group{
		brokers:   brokers,
		topics:    topics,
		groupName: groupName,
	}
}

func (g *group) Start() error {
	log.Info(fmt.Sprintf("Starting a Consumer Group %s", g.groupName))

	config := sarama.NewConfig()
	config.Version = sarama.V1_1_1_0
	config.Consumer.Return.Errors = true

	client, err := sarama.NewClient(g.brokers, config)
	if err != nil {
		log.Fatal("Failed to create kafka client: ", err)
	}
	defer func() { _ = client.Close() }()

	consumer := consumer{
		ready: make(chan bool, 0),
	}

	ctx := context.Background()
	consumerGroup, err := sarama.NewConsumerGroupFromClient(g.groupName, client)
	if err != nil {
		return consumerGroupException.GroupCreation{GroupName: g.groupName, Reasons: []string{err.Error()}}
	}

	go func() {
		for {
			err := consumerGroup.Consume(ctx, g.topics, &consumer)
			if err != nil {
				log.Fatal(consumerGroupException.Consumption{Reasons: []string{err.Error()}}.Error())
			}
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Info(fmt.Sprintf("Consumer Group %s up and running", g.groupName))

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm // Await a sigterm signal before safely closing the consumer

	err = client.Close()
	if err != nil {
		log.Fatal("error closing %s group client: ", err.Error())
	}

	return nil
}
