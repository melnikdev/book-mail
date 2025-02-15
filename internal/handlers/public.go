package handlers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	db *redis.Client
}

func NewHandler(db *redis.Client) *Handler {
	return &Handler{db}
}

func (h *Handler) HelloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (h *Handler) ListUsers(c *fiber.Ctx) error {
	val, err := h.db.LRange(context.Background(), "users:ids", 0, 100).Result()
	if err == redis.Nil {
		fmt.Println("value not found")
	} else if err != nil {
		fmt.Printf("failed to get value, error: %v\n", err)
	}
	fmt.Printf("List of users: %v\n", val)
	return c.SendString("List of users: " + fmt.Sprintf("%v", val))
}
