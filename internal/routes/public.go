package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melnikdev/book-mail/internal/handlers"
	"github.com/redis/go-redis/v9"
)

func PublicRoutes(a *fiber.App, db *redis.Client) {
	handler := handlers.NewHandler(db)

	// Create route group.
	route := a.Group("/api/v1")

	route.Get("/", handler.HelloWorld)

	route.Get("/users", handler.ListUsers)
}
