package memberroute

import (
	"github.com/gofiber/fiber/v3"
	"rishik.com/controller/membercontroller"
	"rishik.com/middleware"
)

func SetupMemberRoutes(app fiber.Router) {
	member := app.Group("/member")
	member.Post("/workspace/invite/:inviteCode/join", middleware.IsAuthenticated(), membercontroller.InviteController)
}
