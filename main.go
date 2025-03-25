package main

import "github.com/gofiber/fiber/v2"

func main() {

	// Create fiber instance
	app := fiber.New(fiber.Config{
		AppName: "Library API",
	})

	// Start server on port 3000
	app.Listen(":3000")
}
