package main

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// jwtSecret is loaded from JWT_SECRET environment variable in main.go init()

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Auth handlers
func Register(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	user.Password = string(hashedPassword)

	if err := DB.Create(user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User created"})
}

func Login(c *fiber.Ctx) error {
	input := struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user User
	if err := DB.Where("user_name = ?", input.UserName).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Update last connection
	now := time.Now()
	user.LastConnectionDate = &now
	DB.Save(&user)

	// Generate JWT
	claims := Claims{
		UserID: user.UserId,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{"token": tokenString})
}

// Middleware to check JWT
func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Missing token"})
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	c.Locals("user_id", claims.UserID)
	c.Locals("role", claims.Role)
	return c.Next()
}

// Role check middleware
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		if role != requiredRole && role != "super_admin" {
			return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
		}
		return c.Next()
	}
}

// User CRUD
func GetUsers(c *fiber.Ctx) error {
	var users []User
	if err := DB.Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	return c.JSON(users)
}

func GetUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	if role != "admin" && role != "super_admin" && userID != id {
		return c.Status(403).JSON(fiber.Map{"error": "Can only view own profile"})
	}

	var user User
	if err := DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	if role != "admin" && role != "super_admin" && userID != id {
		return c.Status(403).JSON(fiber.Map{"error": "Can only update own profile"})
	}

	var user User
	if err := DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	updateData := new(User)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Prevent role update unless super_admin
	if role != "super_admin" {
		updateData.Role = user.Role
	}

	if err := DB.Model(&user).Updates(updateData).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}
	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	role := c.Locals("role").(string)

	if role != "super_admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Only super_admin can delete users"})
	}

	if err := DB.Delete(&User{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}
	return c.JSON(fiber.Map{"message": "User deleted"})
}

// Tama CRUD (similar pattern for other tables)
func GetTamas(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	var tamas []Tama
	query := DB.Preload("User").Preload("TamaStats")
	if role != "admin" && role != "super_admin" {
		query = query.Where("user_id = ?", userID)
	}
	if err := query.Find(&tamas).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch tamas"})
	}
	return c.JSON(tamas)
}

func GetTama(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	var tama Tama
	if err := DB.Preload("User").Preload("TamaStats").First(&tama, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Tama not found"})
	}

	if role != "admin" && role != "super_admin" && tama.UserId != userID {
		return c.Status(403).JSON(fiber.Map{"error": "Can only view own tamas"})
	}

	return c.JSON(tama)
}

func CreateTama(c *fiber.Ctx) error {
	tama := new(Tama)
	if err := c.BodyParser(tama); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	userID := c.Locals("user_id").(int)
	tama.UserId = userID

	if err := DB.Create(tama).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create tama"})
	}
	return c.Status(201).JSON(tama)
}

func UpdateTama(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	var tama Tama
	if err := DB.First(&tama, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Tama not found"})
	}

	if role != "admin" && role != "super_admin" && tama.UserId != userID {
		return c.Status(403).JSON(fiber.Map{"error": "Can only update own tamas"})
	}

	updateData := new(Tama)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := DB.Model(&tama).Updates(updateData).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update tama"})
	}
	return c.JSON(tama)
}

func DeleteTama(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	var tama Tama
	if err := DB.First(&tama, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Tama not found"})
	}

	if role != "super_admin" && (role != "admin" || tama.UserId != userID) {
		return c.Status(403).JSON(fiber.Map{"error": "Insufficient permissions"})
	}

	if err := DB.Delete(&tama, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete tama"})
	}
	return c.JSON(fiber.Map{"message": "Tama deleted"})
}

// Similar handlers for other tables: TamaStats, Friends, Sponsor, Race, Sickness, Trait, Bonus, Malus, Event, LifeChoice
// For brevity, I'll add placeholders

func GetRaces(c *fiber.Ctx) error {
	var races []Race
	if err := DB.Find(&races).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch races"})
	}
	return c.JSON(races)
}

// Add more handlers as needed...
