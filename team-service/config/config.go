package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	ServerPort     string
	UserServiceURL string
	RedisAddr      string
	RedisPassword  string
	RedisDB        int
	KafkaBrokers   []string
	KafkaTopic     string
}

func LoadConfig() (*Config, error) {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}
	return &Config{
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		ServerPort:     os.Getenv("SERVER_PORT"),
		UserServiceURL: os.Getenv("USER_SERVICE_URL"),
		RedisAddr:      os.Getenv("REDIS_ADDR"),
		RedisPassword:  os.Getenv("REDIS_PASSWORD"),
		RedisDB:        redisDB,
		KafkaBrokers:   strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		KafkaTopic:     os.Getenv("KAFKA_TOPIC"),
	}, nil
}
