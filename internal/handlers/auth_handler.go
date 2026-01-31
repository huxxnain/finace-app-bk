package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/huxxnainali/finance-app/internal/auth"
	"github.com/huxxnainali/finance-app/internal/config"
	"github.com/huxxnainali/finance-app/internal/models"
	"github.com/huxxnainali/finance-app/internal/services"
)

type AuthHandler struct {
	userService *services.UserService
	config      *config.Config
}

func NewAuthHandler(userService *services.UserService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		config:      cfg,
	}
}

// SignUp handles user registration
// POST /auth/signup
func (ah *AuthHandler) SignUp(c *fiber.Ctx) error {
	var req models.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email and password are required",
		})
	}

	user, err := ah.userService.SignUp(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID.Hex(), ah.config.JWTSecret, ah.config.JWTExpiryHours)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.AuthResponse{
		Token: token,
	})
}

// Login handles user authentication
// POST /auth/login
func (ah *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email and password are required",
		})
	}

	user, err := ah.userService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID.Hex(), ah.config.JWTSecret, ah.config.JWTExpiryHours)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.AuthResponse{
		Token: token,
	})
}
