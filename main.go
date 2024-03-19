package main

import (
	"log"
	"time"

	"example.com/server/models"
	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
)

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

	// Connect to database
	err := connectDb()
	if err != nil {
		panic(err)
	}

	// Perform auto migration to create table if not present
	db.AutoMigrate(&models.User{})

	// Create a slice with few users
	users := []models.User{
		models.User{Id: 1, Name: "Elon Musk"},
		models.User{Id: 2, Name: "Sam Altman"},
	}
	// Insert users into database if not already exists
	for _, user := range users {
		db.FirstOrCreate(&user, user)
	}

	app := fiber.New() // Notice similarity with express.js => const app = express()

	jwtMiddleware := NewAuthMiddleware(SECRET_KEY)

	// Return simple string
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	// Get request returning JSON
	app.Get("/users", func(c *fiber.Ctx) error {
		var allUsers []models.User
		rows := db.Table("users").Select("id", "name")
		err := rows.Scan(&allUsers).Error
		if err != nil {
			return c.Status(500).SendString("Error fetching users from database")
		}
		return c.JSON(allUsers)
	})

	// Post requesting accepting Json data
	app.Post("/users", jwtMiddleware, func(c *fiber.Ctx) error {
		userDetails := new(models.User)
		err := c.BodyParser(userDetails)
		if err != nil {
			return c.Status(400).SendString("Invalid JSON")
		}
		db.Create(&userDetails)
		return c.JSON(userDetails)
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
