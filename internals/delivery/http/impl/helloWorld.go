package impl

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/triasbrata/adios/internals/entities"
	"github.com/triasbrata/adios/pkgs/instrumentation"
)

func (h *httpHandler) CurrentWether(c *fiber.Ctx) error {
	ctx, span := instrumentation.Tracer.Start(c.UserContext(), "internals:delivery:http:impl:CurrentWether")
	log.Printf("http trace=%s span=%s", span.SpanContext().TraceID().String(), span.SpanContext().SpanID())
	c.SetUserContext(ctx)
	defer span.End()
	latRaw := c.Query("lat", "1")
	lngRaw := c.Query("lng", "1")
	lat, err := strconv.ParseFloat(latRaw, 32)
	if err != nil {
		return err
	}
	lng, err := strconv.ParseFloat(lngRaw, 32)
	if err != nil {
		return err
	}
	data, err := h.service.FetchCurrentWeather(ctx, entities.FetchCurrentWeatherParam{
		Latitude:  float32(lat),
		Longitude: float32(lng),
	})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(data)
}
