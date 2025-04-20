package usercontroller

import (
	"github.com/gofiber/fiber/v3"
	"rishik.com/services/user"
)

func GetCurrentController(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	user, err := user.GetCurrentUsersService(userID.(string))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User fetched Sucessfully",
		"user":    user,
	})
}
