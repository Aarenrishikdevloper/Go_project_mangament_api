package userroute

import (
	"github.com/gofiber/fiber/v3"
	"rishik.com/asynchandler"
	"rishik.com/controller/usercontroller"
	"rishik.com/middleware"
)

func SetupUserRoutes(app fiber.Router) {
	user := app.Group("/users")
	user.Get("/current", middleware.IsAuthenticated(), asynchandler.AsyncHandler(usercontroller.GetCurrentController))
}
