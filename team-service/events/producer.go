package events

import (
	"log"
	"team-service/config"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(cfg *config.Config) (*KafkaProducer, error) {
	producer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, sarama.NewConfig())
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: producer}, nil // Trả về *KafkaProducer
}

func (p *KafkaProducer) PublishMessage(topic string, key string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message published to topic %s [partition: %d, offset: %d]", topic, partition, offset)
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}
