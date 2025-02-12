package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/melnikdev/book-mail/config"
	"github.com/melnikdev/book-mail/internal/services/kafka"
	"github.com/melnikdev/book-mail/internal/services/mail"
)

func main() {
	var user mail.User

	config := config.MustLoad()

	r := kafka.New(config).GetReader()
	m := mail.New(config)

	for {
		msg, err := r.ReadMessage(context.Background())

		if err != nil {
			break
		}

		if err = json.Unmarshal(msg.Value, &user); err != nil {
			log.Printf("Error unmarshaling JSON: %v\n", err)
			continue
		}

		if err := m.Send(&user); err != nil {
			log.Fatal("Failed mail send:", err)
			panic(err)
		}

	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

}
