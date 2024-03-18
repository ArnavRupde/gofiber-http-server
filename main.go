package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	internal_id string // Small case for unexported field
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {

	// Create a slice with users
	users := []User{
		User{Id: 1, Name: "Elon Musk", internal_id: "123"},
		User{Id: 2, Name: "Sam Altman", internal_id: "456"},
	}

	app := fiber.New() // Notice similarity with express.js => const app = express()

	// Return simple string
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	// Get request returning JSON
	app.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(users)
	})

	// Post requesting accepting Json data
	app.Post("/users", func(c *fiber.Ctx) error {
		user := new(User)
		err := c.BodyParser(user)
		if err != nil {
			return c.Status(400).SendString("Invalid JSON")
		}
		users = append(users, *user)
		return c.JSON(users)
	})

	log.Fatal(app.Listen(":3001"))
}
