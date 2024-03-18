package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
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

var SECRET_KEY = "abc123"

// Middleware for JWT
func NewAuthMiddleware(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
	})
}

func main() {

	// Create a slice with users
	users := []User{
		User{Id: 1, Name: "Elon Musk", internal_id: "123"},
		User{Id: 2, Name: "Sam Altman", internal_id: "456"},
	}

	app := fiber.New() // Notice similarity with express.js => const app = express()

	jwtMiddleware := NewAuthMiddleware(SECRET_KEY)

	// Return simple string
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	// Get request returning JSON
	app.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(users)
	})

	// Post requesting accepting Json data
	app.Post("/users", jwtMiddleware, func(c *fiber.Ctx) error {
		user := new(User)
		err := c.BodyParser(user)
		if err != nil {
			return c.Status(400).SendString("Invalid JSON")
		}
		users = append(users, *user)
		return c.JSON(users)
	})

	// Login route
	app.Post("/login", func(c *fiber.Ctx) error {
		loginCredentials := new(Login)
		err := c.BodyParser(loginCredentials)
		if err != nil && loginCredentials.Username != "" && loginCredentials.Password != "" {
			return c.Status(400).SendString("Invalid input JSON")
		}

		username := "admin"
		password := "admin"
		if loginCredentials.Username == username && loginCredentials.Password == password {
			// Set up content for JWT
			claims := jwt.MapClaims{
				"username": username,
				"exp":      time.Now().Add(time.Hour * 24).Unix(),
			}
			// Create token
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			//Encode token
			encodedToken, err := token.SignedString([]byte(SECRET_KEY))
			if err != nil {
				return c.SendStatus(500)
			}
			return c.JSON(fiber.Map{"token": encodedToken})
		} else {
			return c.SendStatus(401)
		}
	})

	log.Fatal(app.Listen(":3001"))
}
