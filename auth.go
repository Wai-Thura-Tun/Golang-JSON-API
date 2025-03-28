package main

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AuthHandlers(route fiber.Router, db *gorm.DB) {
	route.Post("/register", func(c *fiber.Ctx) error {
		authUser := &User{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}

		if authUser.Username == "" || authUser.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username and password required",
			})
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(authUser.Password), bcrypt.DefaultCost)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		authUser.Password = string(hashed)
		db.Create(authUser)

		token, err := GenerateToken(authUser)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			HTTPOnly: !c.IsFromLocal(),
			Secure:   !c.IsFromLocal(),
			MaxAge:   3600 * 24 * 7,
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": token,
		})
	})

	route.Post("/login", func(c *fiber.Ctx) error {
		dbUser := new(User)
		authUser := &User{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}

		db.Where("username = ?", authUser.Username).First(&dbUser)

		if dbUser.ID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "User not found.",
			})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(authUser.Password)); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid Password",
			})
		}

		token, err := GenerateToken(authUser)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			HTTPOnly: !c.IsFromLocal(),
			Secure:   !c.IsFromLocal(),
			MaxAge:   3600 * 24 * 7,
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": token,
		})
	})
}
