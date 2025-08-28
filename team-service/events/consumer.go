package events

import (
	"log"
	"team-service/config"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
}

func NewKafkaConsumer(cfg *config.Config) (*KafkaConsumer, error) {
	consumer, err := sarama.NewConsumer(cfg.KafkaBrokers, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{consumer: consumer}, nil
}

func (c *KafkaConsumer) Subscribe(topic string, handler func(message *sarama.ConsumerMessage)) error {
	partitionConsumer, err := c.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}

	go func() {
		for msg := range partitionConsumer.Messages() {
			log.Printf("Received message on topic %s: %s", msg.Topic, string(msg.Value))
			handler(msg)
		}
	}()

	return nil
}

func (c *KafkaConsumer) Close() error {
	return c.consumer.Close()
}

