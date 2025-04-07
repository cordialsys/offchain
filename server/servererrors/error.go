package servererrors

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

func BadRequestf(c *fiber.Ctx, msg string, args ...interface{}) error {
	return NewErrorf(c, fiber.StatusBadRequest, msg, args...)
}

func Unauthorizedf(c *fiber.Ctx, msg string, args ...interface{}) error {
	return NewErrorf(c, fiber.StatusUnauthorized, msg, args...)
}

// Maps to Aborted
func Conflictf(c *fiber.Ctx, msg string, args ...interface{}) error {
	return NewErrorf(c, fiber.StatusConflict, msg, args...)
}
