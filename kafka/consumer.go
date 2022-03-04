package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

type MsgCallback chan []byte

type consumer struct {
	client sarama.ConsumerGroup
	topics map[string]MsgCallback
}

type Consumer interface {
	Serve(ctx context.Context)
	Close()
}

func NewConsumer(brokers []string, topics map[string]MsgCallback) (Consumer, error) {
	config, err := consumerConfig()
	if err != nil {
		return nil, err
	}

	client, err := sarama.NewConsumerGroup(brokers, "t-pakhuis-consumer", config)
	if err != nil {
		return nil, err
	}

	consumer := &consumer{
		client: client,
		topics: topics,
	}
	return consumer, nil
}

func consumerConfig() (*sarama.Config, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	version, err := sarama.ParseKafkaVersion("2.1.1")
	if err != nil {
		return nil, err
	}
	config.Version = version
	return config, nil
}

func (c *consumer) Serve(ctx context.Context) {
	topics := keysFromMap(c.topics)
	for {
		if err := c.client.Consume(ctx, topics, c); err != nil {
			log.Panicf("Error from consumer: %v", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

func (c *consumer) Close() {
	c.client.Close()
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	log.Info("Kafka consumer is setup")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		c.topics[message.Topic] <- message.Value
		session.MarkMessage(message, "")
	}
	return nil
}

func keysFromMap(m map[string]MsgCallback) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
