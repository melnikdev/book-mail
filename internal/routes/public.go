package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melnikdev/book-mail/internal/handlers"
	"github.com/redis/go-redis/v9"
)

func PublicRoutes(r fiber.Router, db *redis.Client) {
	handler := handlers.NewHandler(db)

	r.Get("/", handler.HelloWorld)
	r.Get("/users", handler.ListUsers)
}
