package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melnikdev/book-mail/internal/handlers"
)

func PublicRoutes(a *fiber.App) {
	// Create route group.
	route := a.Group("/api/v1")

	route.Get("/", handlers.HelloWorld)

	route.Get("/users", handlers.ListUsers)
}
