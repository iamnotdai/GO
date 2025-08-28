package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"team-service/cache"
	"team-service/config"
	"team-service/events"

	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
)

type TeamConsumer struct {
	kafkaConsumer *events.KafkaConsumer
	redisClient   *redis.Client
}

func NewTeamConsumer(cfg *config.Config) (*TeamConsumer, error) {
	kafkaConsumer, err := events.NewKafkaConsumer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %v", err)
	}
	redisClient := cache.NewRedisClient(cfg)
	return &TeamConsumer{
		kafkaConsumer: kafkaConsumer,
		redisClient:   redisClient,
	}, nil
}

func (tc *TeamConsumer) Start(ctx context.Context) error {
	err := tc.kafkaConsumer.Subscribe(events.TopicTeamCreated, func(message *sarama.ConsumerMessage) {
		// Parse dữ liệu từ Kafka message
		var teamEvent struct {
			TeamID   string `json:"team_id"`
			TeamName string `json:"team_name"`
		}
		if err := json.Unmarshal(message.Value, &teamEvent); err != nil {
			log.Printf("failed to unmarshal team event: %v", err)
			return
		}

		// Lưu vào Redis
		key := fmt.Sprintf("team:%s", teamEvent.TeamID)
		value, err := json.Marshal(teamEvent)
		if err != nil {
			log.Printf("failed to marshal team event for redis: %v", err)
			return
		}

		err = tc.redisClient.Set(ctx, key, value, 0).Err()
		if err != nil {
			log.Printf("failed to save to redis: %v", err)
			return
		}
		log.Printf("Saved team %s to redis", teamEvent.TeamID)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %v", events.TopicTeamCreated, err)
	}

	// Giữ consumer chạy
	go func() {
		<-ctx.Done()
		tc.kafkaConsumer.Close()
		tc.redisClient.Close()
	}()

	return nil
}

func (tc *TeamConsumer) Close() error {
	if err := tc.kafkaConsumer.Close(); err != nil {
		return err
	}
	return tc.redisClient.Close()
}
