package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/melnikdev/book-mail/config"
	"github.com/melnikdev/book-mail/database"
	"github.com/melnikdev/book-mail/internal/services/kafka"
	"github.com/melnikdev/book-mail/internal/services/mail"
	"github.com/melnikdev/book-mail/server"
)

func main() {
	config := config.MustLoad()
	cu := make(chan mail.User)

	db, err := database.NewClient(context.Background())
	if err != nil {
		panic(err)
	}

	server := server.NewFiberServer(config, db)

	r := kafka.New(config).GetReader()
	m := mail.New(config)

	go kafka.Consumer(r, db, cu)
	go m.ListenEvent(cu)

	go server.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c // Block until a signal is received

	_ = server.Shutdown(r)
}
