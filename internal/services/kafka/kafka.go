package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/melnikdev/book-mail/config"
	"github.com/melnikdev/book-mail/internal/services/mail"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type kafkaBroker struct {
	config *config.Config
}

func New(config *config.Config) *kafkaBroker {
	return &kafkaBroker{config: config}
}

func (b *kafkaBroker) GetReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{b.config.Kafka.Broker},
		GroupID:  b.config.Kafka.GroupID,
		Topic:    b.config.Kafka.Topic,
		MaxBytes: b.config.Kafka.MaxBytes, // 10MB
	})
}

func Consumer(k *kafka.Reader, db *redis.Client, ch chan mail.User) {
	var user mail.User

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			msg, err := k.ReadMessage(context.Background())

			if err != nil {
				continue
			}

			if err = json.Unmarshal(msg.Value, &user); err != nil {
				log.Printf("Error unmarshaling JSON: %v\n", err)
				continue
			}
			ch <- user

			if err := db.LPush(context.Background(), "users:ids", user.UserId).Err(); err != nil {
				fmt.Printf("failed to set data, error: %s", err.Error())
			}
		}
	}
}
