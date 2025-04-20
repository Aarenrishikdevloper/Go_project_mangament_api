package authroute

import (
	"github.com/gofiber/fiber/v3"

	"rishik.com/asynchandler"
	"rishik.com/controller/authcontroller"
)

func SetupAuthRoutes(app fiber.Router) {

	auth := app.Group("/auth")
	auth.Post("/register", asynchandler.AsyncHandler(authcontroller.RegisterUserController))

	auth.Post("/login", asynchandler.AsyncHandler(authcontroller.LoginUserController))
	auth.Post("/logout", asynchandler.AsyncHandler(authcontroller.LogoutController))
	auth.Get("/test", func(c fiber.Ctx) error {
		return c.SendString("Auth routes are working!")
	})
}
