package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/melnikdev/book-mail/config"
	"github.com/melnikdev/book-mail/internal/services/kafka"
	"github.com/melnikdev/book-mail/internal/services/mail"
)

func main() {
	var user mail.User

	config := config.MustLoad()

	r := kafka.New(config).GetReader()
	m := mail.New(config)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			msg, err := r.ReadMessage(context.Background())

			if err != nil {
				continue
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
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

}
