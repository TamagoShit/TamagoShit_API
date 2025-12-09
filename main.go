package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	JWTSecret []byte
)

func init() {
	// Try to load environment variables from Render's secrets file
	if err := godotenv.Load("/etc/secrets/.env"); err != nil {
		log.Printf("Warning: Could not load .env file from /etc/secrets/.env: %v", err)
		log.Println("Falling back to system environment variables")
	}

	// Load JWT secret from environment or use default
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("Warning: JWT_SECRET not set, using default secret (not recommended for production)")
		secret = "your-default-jwt-secret-change-this-in-production"
	}
	JWTSecret = []byte(secret)
}

func main() {
	// Get DATABASE_URL from environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	err = DB.AutoMigrate(&User{}, &Tama{}, &TamaStats{}, &Friend{}, &Sponsor{}, &Race{}, &Sickness{}, &Trait{}, &Bonus{}, &Malus{}, &Event{}, &LifeChoice{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	app := fiber.New()

	// Add CORS middleware
	app.Use(cors.New())

	// Public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("TamagoShit API")
	})

	app.Post("/register", Register)
	app.Post("/login", Login)

	// Protected routes
	api := app.Group("/api", AuthMiddleware)

	// User routes
	api.Get("/users", RequireRole("admin"), GetUsers)
	api.Get("/users/:id", GetUser)
	api.Put("/users/:id", UpdateUser)
	api.Delete("/users/:id", RequireRole("super_admin"), DeleteUser)

	// Tama routes
	api.Get("/tamas", GetTamas)
	api.Get("/tamas/:id", GetTama)
	api.Post("/tamas", CreateTama)
	api.Put("/tamas/:id", UpdateTama)
	api.Delete("/tamas/:id", DeleteTama)

	// Race routes (public read)
	app.Get("/api/races", GetRaces)

	// Add more routes for other tables...

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
