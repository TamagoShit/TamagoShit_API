package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Database connection
	dsn := "general_user:password@tcp(localhost:3306)/tamagoshit_db?charset=utf8mb4&parseTime=True&loc=Local" // Update with actual credentials
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
