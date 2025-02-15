package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/melnikdev/book-mail/config"
	"github.com/melnikdev/book-mail/database"
	"github.com/melnikdev/book-mail/internal/routes"
	"github.com/melnikdev/book-mail/internal/services/kafka"
	"github.com/melnikdev/book-mail/internal/services/mail"
)

func main() {
	config := config.MustLoad()
	cu := make(chan mail.User)

	db, err := database.NewClient(context.Background())

	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(healthcheck.New())

	routes.PublicRoutes(app, db)

	r := kafka.New(config).GetReader()
	m := mail.New(config)

	go kafka.Consumer(r, db, cu)
	go m.ListenEvent(cu)

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c // Block until a signal is received

	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

	fmt.Println("Fiber was successful shutdown.")
}
