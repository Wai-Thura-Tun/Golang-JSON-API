package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Download(route fiber.Router, db *gorm.DB) {
	route.Get("/", func(c *fiber.Ctx) error {
		format := c.Query("format", "json")
		books := new([]Book)

		if err := db.Find(books).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book Not Found",
			})
		}

		var fileName string
		switch format {
		case "json":
			fileName = "books.json"
			file, err := os.Create(fileName)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create the JSON file.",
				})
			}

			defer file.Close()

			// Write book data to the JSON File
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", " ")
			if err := encoder.Encode(books); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to write the JSON file",
				})
			}

		case "csv":
			fileName = "books.csv"
			file, err := os.Create(fileName)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create the JSON file.",
				})
			}

			defer file.Close()

			// Write book data to CSV File
			writer := csv.NewWriter(file)

			// CSV Header
			writer.Write([]string{"ID", "Title", "Status", "Author", "Year"})
			for _, book := range *books {
				writer.Write([]string{
					fmt.Sprintf("%d", book.ID),
					book.Title,
					string(book.Status),
					book.Author,
					fmt.Sprintf("%d", book.Year),
				})
			}

			writer.Flush()
			if err := writer.Error(); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Failed to write CSV file",
				})
			}

		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid File Format. User json or csv instead",
			})
		}

		defer os.Remove(fileName)
		return c.Download(fileName)
	})
}
