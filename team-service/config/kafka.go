package config

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

var KafkaProducer sarama.SyncProducer

func InitKafka(cfg *Config) (sarama.SyncProducer, error) {
	if cfg == nil || len(cfg.KafkaBrokers) == 0 {
		return nil, fmt.Errorf("missing kafka brokers in config")
	}

	// cấu hình sarama tốt hơn là truyền nil
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.RequiredAcks = sarama.WaitForAll
	saramaCfg.Producer.Retry.Max = 5
	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.Timeout = 5 * time.Second

	producer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, saramaCfg)
	if err != nil {
		log.Printf("Failed to connect Kafka brokers %v: %v", cfg.KafkaBrokers, err)
		return nil, err
	}

	KafkaProducer = producer
	log.Println("Connected to Kafka at", cfg.KafkaBrokers)
	return producer, nil
}

// CloseKafka đóng Kafka producer
func CloseKafka() {
	if KafkaProducer != nil {
		err := KafkaProducer.Close()
		if err != nil {
			log.Printf("Failed to close Kafka producer: %v", err)
		}
	}
}
