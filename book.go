package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BookHandlers(route fiber.Router, db *gorm.DB) {
	route.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Get All")
	})

	route.Get("/:id", func(c *fiber.Ctx) error {
		return c.SendString("Get By Id")
	})

	route.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("Create Book")
	})

	route.Put("/:id", func(c *fiber.Ctx) error {
		return c.SendString("Update Book")
	})

	route.Delete("/:id", func(c *fiber.Ctx) error {
		return c.SendString("Delete Book")
	})
}
