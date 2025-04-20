package authcontroller

import (
	"fmt"
	"log"

	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"

	"go.mongodb.org/mongo-driver/mongo"
	"rishik.com/services/auth"
)

func RegisterUserController(c fiber.Ctx) error {
	var body auth.RegisterRequest
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if body.Email == "" || body.Name == "" || body.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "All fields are Required")
	}
	user, err := auth.RegisterUserService(body)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fiber.NewError(fiber.StatusBadRequest, "Email already exists")

		}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	var jwtsecret = []byte(os.Getenv("SESSION_SECRET"))
	tokenString, err := token.SignedString(jwtsecret)
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Could not generate token")

	}
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: false,
		Secure:   false,
		SameSite: "None",
	})
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User Created sucessfully",
		"user":    user,
	})
}

func LoginUserController(c fiber.Ctx) error {
	var body auth.LoginRequest
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if body.Email == "" || body.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "All fields are Required")
	}
	user, err := auth.LoginService(body)
	if err != nil {
		log.Fatalln("faile to clear roles: %w", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Something went Wrong")

	}
	var jwtsecret = []byte(os.Getenv("SESSION_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	fmt.Print(jwtsecret)
	tokenString, err := token.SignedString(jwtsecret)

	if err != nil {
		log.Printf("failed to generate token: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Could not generate token")

	}
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "lax",
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "User login Sucessfully ",
		"user":      user,
		tokenString: tokenString,
	})
}
func LogoutController(c fiber.Ctx) error {
	c.ClearCookie("auth_token")
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User logged out successfully",
	})
	return nil
}
