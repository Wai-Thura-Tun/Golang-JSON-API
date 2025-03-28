package main

import "github.com/gofiber/fiber/v2"

func main() {

	// Initalize database
	db := InitializeDB()

	// Create fiber instance
	app := fiber.New(fiber.Config{
		AppName: "Library API",
	})

	api := app.Group("/api")

	// Define auth routhes. Those will be public
	AuthHandlers(api.Group("/auth"), db)

	// Verify the JWT
	protected := api.Group("/book", AuthMiddleware(db))

	// Define book routhes
	BookHandlers(protected, db)

	// Start server on port 3000
	app.Listen(":3000")
}
