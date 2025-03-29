package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BookHandlers(route fiber.Router, db *gorm.DB) {
	route.Get("/", func(c *fiber.Ctx) error {
		title := c.Query("title")
		status := c.Query("status")
		author := c.Query("author")
		year := c.QueryInt("year")

		userId := int(c.Locals("userId").(float64))
		books := new([]Book)

		// Init Query
		query := db.Where("user_id = ?", userId)

		if title != "" {
			query.Where("title LIKE ?", "%"+title+"%")
		}

		if status != "" {
			query.Where("status = ?", status)
		}

		if author != "" {
			query.Where("author = ?", author)
		}

		if year != 0 {
			query.Where("year = ?", year)
		}

		if err := query.Find(books).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book Not Found",
			})
		}
		return c.Status(fiber.StatusOK).JSON(books)
	})

	route.Get("/:id", func(c *fiber.Ctx) error {
		bookId := c.Params("id")
		userId := int(c.Locals("userId").(float64))
		book := new(Book)
		if err := db.Where("id = ? AND user_id = ?", bookId, userId).First(book).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book Not Found",
			})
		}
		return c.Status(fiber.StatusOK).JSON(book)
	})

	route.Post("/", func(c *fiber.Ctx) error {
		book := new(Book)
		book.UserID = uint(c.Locals("userId").(float64))

		if err := c.BodyParser(book); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := db.Create(book).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(book)
	})

	route.Put("/:id", func(c *fiber.Ctx) error {
		bookId, _ := c.ParamsInt("id")
		userId := int(c.Locals("userId").(float64))
		book := new(Book)

		if err := db.Where("id = ? AND user_id = ?", bookId, userId).First(book).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book Not Found",
			})
		}

		if err := c.BodyParser(book); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := db.Save(book).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(book)
	})

	route.Delete("/:id", func(c *fiber.Ctx) error {

		bookId, _ := c.ParamsInt("id")
		userId := int(c.Locals("userId").(float64))
		book := new(Book)

		if err := db.Where("id = ? AND user_id = ?", bookId, userId).First(book).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Book Not Found",
			})
		}

		if err := db.Delete(book).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
