package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//Get the token either from the header or cookie
		cookieToken := c.Cookies("jwt")
		var tokenString string
		if cookieToken != "" {
			tokenString = cookieToken
			log.Print("Using cookie token...")
		} else {
			// Get the auth header
			authHeader := c.Get("Authorization")

			if authHeader == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Unauthorized",
				})
			}

			// Split the header
			tokenParts := strings.Split(authHeader, "")

			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Unauthorized",
				})
			}
			tokenString = tokenParts[1]
			log.Print("Using token from auth header...")
		}

		// Parse Token
		secret := []byte("super-secret-key")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.ClearCookie()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Get user from the token
		userId := token.Claims.(jwt.MapClaims)["userId"]

		// Check if user exists in the db, if not clear the cookie
		if err := db.Model(&User{}).Where("id = ?", userId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			c.ClearCookie()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Set the userId in the locals
		c.Locals("userId", userId)

		return c.Next()
	}
}
