package main

import (
	"log"
	"time"

	"example.com/server/auth"
	"example.com/server/models"
	"example.com/server/services"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/contrib/websocket"
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
	db.AutoMigrate(&models.Wallet{})

	// Create a slice with few users
	users := []models.User{
		models.User{Id: 1, Name: "Elon Musk"},
		models.User{Id: 2, Name: "Sam Altman"},
	}
	// Insert users into database if not already exists
	for _, user := range users {
		db.FirstOrCreate(&user, user)
	}

	// Create a wallet for user with id 1
	db.Create(
		&models.Wallet{
			UserId:  1,
			Balance: decimal.NewFromFloat(1000.00),
		},
	)
	db.Create(
		&models.Wallet{
			UserId:  2,
			Balance: decimal.NewFromFloat(1000.00),
		},
	)

	paymentService := services.NewPaymentService(db)

	app := fiber.New() // Notice similarity with express.js => const app = express()

	jwtMiddleware := NewAuthMiddleware(SECRET_KEY)

	// Return simple string
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		for {
			// Read message from websocket connection
			mtype, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			log.Printf("Message received: %s", msg)
			// Write message back to websocket connection
			err = c.WriteMessage(mtype, msg)
			if err != nil {
				break
			}
		}
		log.Printf("Error: %v", err)
	}))

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

	// Signup route
	app.Post("/signup", func(c *fiber.Ctx) error {
		signupCredentials := new(Login)
		err := c.BodyParser(signupCredentials)
		if err != nil && signupCredentials.Username != "" && signupCredentials.Password != "" {
			return c.Status(400).SendString("Invalid input JSON")
		}
		// Check if user already exists
		var users []models.User
		db.Where("username = ?", signupCredentials.Username).Find(&users)
		if len(users) > 0 {
			return c.Status(400).SendString("User already exists")
		}
		// Create user
		hashedPassword, err := auth.HashPassword(signupCredentials.Password)
		if err != nil {
			return c.Status(500).SendString("Error creating user")
		}
		newUser := models.User{
			Name:           "New User",
			Username:       signupCredentials.Username,
			HashedPassword: hashedPassword,
		}
		result := db.Create(&newUser)
		if result.Error != nil {
			return c.Status(500).SendString("Error creating user")
		}
		return c.JSON(newUser)
	})

	// Login route
	app.Post("/login", func(c *fiber.Ctx) error {
		loginCredentials := new(Login)
		err := c.BodyParser(loginCredentials)
		if err != nil && loginCredentials.Username != "" && loginCredentials.Password != "" {
			return c.Status(400).SendString("Invalid input JSON")
		}

		username := loginCredentials.Username
		password := loginCredentials.Password
		// Verify that user exists
		var users []models.User
		db.Where("username = ?", username).Find(&users)
		if len(users) == 0 {
			return c.Status(400).SendString("User does not exist")
		}
		// Verify password
		hashedPassword := users[0].HashedPassword
		err = auth.ComaparePassword(hashedPassword, password)
		if err != nil {
			return c.Status(401).SendString("Incorrect password")
		}

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
	})

	app.Post("/transfer_funds", func(c *fiber.Ctx) error {
		transferDetails := new(services.TransferDetails)
		err := c.BodyParser(transferDetails)
		if err != nil || transferDetails.FromUserId <= 0 || transferDetails.ToUserId <= 0 || !transferDetails.Amount.IsPositive() {
			return c.Status(400).SendString("Invalid input JSON")
		}
		err = paymentService.TransferFunds(
			c.Context(),
			transferDetails.FromUserId,
			transferDetails.ToUserId,
			transferDetails.Amount,
		)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.SendString("Funds transferred successfully")
	})

	log.Fatal(app.Listen(":3001"))
}
