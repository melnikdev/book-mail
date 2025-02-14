package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/melnikdev/book-mail/config"
	"github.com/melnikdev/book-mail/internal/services/mail"
	"github.com/segmentio/kafka-go"
)

type KafkaBroker struct {
	Config *config.Config
}

func New(config *config.Config) *KafkaBroker {
	return &KafkaBroker{Config: config}
}

func (b *KafkaBroker) GetReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{b.Config.Kafka.Broker},
		GroupID:  b.Config.Kafka.GroupID,
		Topic:    b.Config.Kafka.Topic,
		MaxBytes: b.Config.Kafka.MaxBytes, // 10MB
	})
}

func Consumer(k *kafka.Reader, m *mail.Mail, ch chan mail.User) {
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

			if err := m.Send(&user); err != nil {
				log.Fatal("Failed mail send:", err)
				panic(err)
			}
		}
	}
}
