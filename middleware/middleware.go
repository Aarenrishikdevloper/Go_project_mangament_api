package middleware

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func IsAuthenticated() fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenString := c.Cookies("auth_token")
		if tokenString == "" {
			fmt.Print(tokenString)
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized. Please Log in")
		}

		// Get the same secret used for signing
		secret := []byte(os.Getenv("SESSION_SECRET"))
		fmt.Print(secret)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil {
			log.Printf("Token validation error: %v", err)
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized. Please Log in")
		}

		if !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["user_id"] == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized. Please Log in")
		}

		c.Locals("userId", claims["user_id"])
		return c.Next()
	}
}
