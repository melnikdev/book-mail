package kafka

import (
	"github.com/melnikdev/book-mail/config"
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
