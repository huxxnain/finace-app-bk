package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/huxxnainali/finance-app/internal/models"
	"github.com/huxxnainali/finance-app/internal/services"
	"github.com/huxxnainali/finance-app/internal/utils"
)

type BudgetHandler struct {
	budgetService *services.BudgetService
}

func NewBudgetHandler(budgetService *services.BudgetService) *BudgetHandler {
	return &BudgetHandler{
		budgetService: budgetService,
	}
}

// GetCurrentBudget retrieves the current month's budget
// GET /budget/current
func (bh *BudgetHandler) GetCurrentBudget(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	year, month := utils.GetCurrentMonthYear()

	budget, err := bh.budgetService.GetOrCreateBudget(c.Context(), userID, year, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	remaining := services.CalculateRemaining(budget.BaseIncome, budget.Expenses)

	return c.Status(fiber.StatusOK).JSON(models.BudgetResponse{
		Year:       budget.Year,
		Month:      budget.Month,
		BaseIncome: budget.BaseIncome,
		Expenses:   budget.Expenses,
		Remaining:  remaining,
	})
}

// GetBudgetByMonth retrieves a specific month's budget
// GET /budget?year=YYYY&month=MM
func (bh *BudgetHandler) GetBudgetByMonth(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "year and month query parameters are required",
		})
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid year parameter",
		})
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid month parameter",
		})
	}

	if month < 1 || month > 12 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "month must be between 1 and 12",
		})
	}

	budget, err := bh.budgetService.GetOrCreateBudget(c.Context(), userID, year, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	remaining := services.CalculateRemaining(budget.BaseIncome, budget.Expenses)

	return c.Status(fiber.StatusOK).JSON(models.BudgetResponse{
		Year:       budget.Year,
		Month:      budget.Month,
		BaseIncome: budget.BaseIncome,
		Expenses:   budget.Expenses,
		Remaining:  remaining,
	})
}

// SetBaseIncome sets or updates the base income for the current month
// POST /budget/base-income
func (bh *BudgetHandler) SetBaseIncome(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	var req models.BaseIncomeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}
	if req.Year <= 0 || req.Month <= 0 || req.Month > 12 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "valid year and month are required",
		})
	}
	if req.Amount < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "amount must be non-negative",
		})
	}

	budget, err := bh.budgetService.SetBaseIncome(c.Context(), userID, req.Year, req.Month, req.Amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	remaining := services.CalculateRemaining(budget.BaseIncome, budget.Expenses)

	return c.Status(fiber.StatusOK).JSON(models.BudgetResponse{
		Year:       budget.Year,
		Month:      budget.Month,
		BaseIncome: budget.BaseIncome,
		Expenses:   budget.Expenses,
		Remaining:  remaining,
	})
}
