package impl

import (
	"github.com/gofiber/fiber/v2"
	"github.com/triasbrata/adios/internals/entities"
)

func (h *httpHandler) HelloWorld(c *fiber.Ctx) error {
	data, err := h.service.FetchHelloWorld(c.UserContext(), entities.FetchHelloWorldParam{})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(data.MapData)
}
