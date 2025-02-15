package handlers

import "github.com/gofiber/fiber/v2"

func HelloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func ListUsers(c *fiber.Ctx) error {
	return c.SendString("List of users")
}
