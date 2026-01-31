package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/huxxnainali/finance-app/internal/auth"
	"github.com/huxxnainali/finance-app/internal/config"
	"github.com/huxxnainali/finance-app/internal/db"
	"github.com/huxxnainali/finance-app/internal/handlers"
	"github.com/huxxnainali/finance-app/internal/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	err := db.Connect(cfg.MongoDBURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Close()

	// Get database instance
	database := db.GetDatabase(cfg.DatabaseName)

	// Initialize services
	userService := services.NewUserService(database)
	budgetService := services.NewBudgetService(database)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, cfg)
	budgetHandler := handlers.NewBudgetHandler(budgetService)
	expenseHandler := handlers.NewExpenseHandler(budgetService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Finance Tracker API v1.0.0",
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type,Authorization",
	}))
	app.Use(logger.New())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Auth routes (no authentication required)
	authGroup := app.Group("/auth")
	authGroup.Post("/signup", authHandler.SignUp)
	authGroup.Post("/login", authHandler.Login)

	// Protected routes (authentication required)
	// Budget routes
	budgetGroup := app.Group("/budget")
	budgetGroup.Use(auth.AuthMiddleware(cfg))
	budgetGroup.Get("/current", budgetHandler.GetCurrentBudget)
	budgetGroup.Get("/", budgetHandler.GetBudgetByMonth)
	budgetGroup.Post("/base-income", budgetHandler.SetBaseIncome)
	budgetGroup.Put("/base-income", budgetHandler.SetBaseIncome)

	// Expense routes
	expenseGroup := app.Group("/expenses")
	expenseGroup.Use(auth.AuthMiddleware(cfg))
	expenseGroup.Post("/", expenseHandler.AddExpense)
	expenseGroup.Put("/:expenseId", expenseHandler.UpdateExpense)
	expenseGroup.Delete("/:expenseId", expenseHandler.DeleteExpense)

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "endpoint not found",
		})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
