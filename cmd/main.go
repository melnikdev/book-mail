package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/melnikdev/book-mail/config"
	"github.com/melnikdev/book-mail/internal/services/kafka"
	"github.com/melnikdev/book-mail/internal/services/mail"
)

func main() {

	config := config.MustLoad()
	cu := make(chan mail.User)
	users := []mail.User{}

	r := kafka.New(config).GetReader()
	m := mail.New(config)

	go kafka.Consumer(r, cu)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(users)
	})

	go func() {
		var user mail.User
		for {
			user = <-cu
			users = append(users, user)
			if err := m.Send(&user); err != nil {
				log.Fatal("Failed mail send:", err)
				panic(err)
			} else {
				fmt.Println("Mail sent to user id: ", user.UserId)
			}
		}
	}()

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
