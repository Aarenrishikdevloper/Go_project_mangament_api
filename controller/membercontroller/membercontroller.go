package membercontroller

import (
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rishik.com/services/member"
)

func InviteController(c fiber.Ctx) error {
	inviteCode := c.Params("inviteCode")
	userId := c.Locals("userId")
	if inviteCode == "" || userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})

	}
	userIdobj, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": " ID of user is Invalid",
		})
	}
	Joindata, err := member.JoinWorkspace(userIdobj, inviteCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to Join workspaces",
		})

	}
	return c.JSON(Joindata)

}
