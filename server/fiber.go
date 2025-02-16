package server

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/melnikdev/book-mail/config"
	"github.com/melnikdev/book-mail/internal/routes"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type fiberServer struct {
	app  *fiber.App
	db   *redis.Client
	conf *config.Config
}

func NewFiberServer(conf *config.Config, db *redis.Client) *fiberServer {
	return &fiberServer{
		app:  fiber.New(),
		db:   db,
		conf: conf,
	}
}

func (s *fiberServer) Start() {
	s.app.Use(recover.New())
	s.app.Use(logger.New())
	s.app.Use(healthcheck.New())

	s.initPublicHttpHandler()

	if err := s.app.Listen(":" + strconv.Itoa(s.conf.Server.Port)); err != nil {
		log.Panic(err)
	}
}

func (s *fiberServer) initPublicHttpHandler() {
	api := s.app.Group("/api/v1/")
	routes.PublicRoutes(api, s.db)
}

func (s *fiberServer) Shutdown(r *kafka.Reader) error {
	fmt.Println("Gracefully shutting down...")

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

	fmt.Println("Fiber was successful shutdown.")
	return s.app.Shutdown()
}
