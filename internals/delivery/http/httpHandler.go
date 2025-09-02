package http

import "github.com/gofiber/fiber/v2"

type Handler interface {
	CurrentWether(c *fiber.Ctx) error
}
