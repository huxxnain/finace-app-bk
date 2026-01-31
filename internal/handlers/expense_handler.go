package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/huxxnainali/finance-app/internal/models"
	"github.com/huxxnainali/finance-app/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpenseHandler struct {
	budgetService *services.BudgetService
}

func NewExpenseHandler(budgetService *services.BudgetService) *ExpenseHandler {
	return &ExpenseHandler{
		budgetService: budgetService,
	}
}

// AddExpense adds a new expense to the current month
// POST /expenses
func (eh *ExpenseHandler) AddExpense(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req models.ExpenseRequest
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

	// Validate input
	if req.Title == "" || req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "title and amount (positive) are required",
		})
	}

	expense := models.Expense{
		ID:        primitive.NewObjectID(),
		Title:     req.Title,
		Amount:    req.Amount,
		CreatedAt: time.Now(),
	}

	budget, err := eh.budgetService.AddExpense(c.Context(), userID, req.Year, req.Month, expense)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	remaining := services.CalculateRemaining(budget.BaseIncome, budget.Expenses)

	return c.Status(fiber.StatusCreated).JSON(models.BudgetResponse{
		Year:       budget.Year,
		Month:      budget.Month,
		BaseIncome: budget.BaseIncome,
		Expenses:   budget.Expenses,
		Remaining:  remaining,
	})
}

// UpdateExpense updates an existing expense
// PUT /expenses/:expenseId
func (eh *ExpenseHandler) UpdateExpense(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	expenseID := c.Params("expenseId")

	var req models.ExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}

	// Validate input
	if req.Title == "" || req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "title and amount (positive) are required",
		})
	}

	updatedExpense := models.Expense{
		Title:  req.Title,
		Amount: req.Amount,
	}

	budget, err := eh.budgetService.UpdateExpense(c.Context(), userID, expenseID, updatedExpense)
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

// DeleteExpense deletes an expense
// DELETE /expenses/:expenseId
func (eh *ExpenseHandler) DeleteExpense(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	expenseID := c.Params("expenseId")

	budget, err := eh.budgetService.DeleteExpense(c.Context(), userID, expenseID)
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
