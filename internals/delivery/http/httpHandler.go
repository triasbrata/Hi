package http

import "github.com/gofiber/fiber/v2"

type Handler interface {
	HelloWorld(c *fiber.Ctx) error
}
