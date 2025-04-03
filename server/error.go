package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func NewErrorf(c *fiber.Ctx, code int, msg string, args ...interface{}) error {
	return c.Status(code).JSON(fiber.Map{
		"message": fmt.Sprintf(msg, args...),
	})
}

func NotFoundf(c *fiber.Ctx, msg string, args ...interface{}) error {
	return NewErrorf(c, fiber.StatusNotFound, msg, args...)
}

func InternalErrorf(c *fiber.Ctx, msg string, args ...interface{}) error {
	return NewErrorf(c, fiber.StatusInternalServerError, msg, args...)
}
