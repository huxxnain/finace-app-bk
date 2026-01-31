package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/huxxnainali/finance-app/internal/models"
	"github.com/huxxnainali/finance-app/internal/services"
)

type FundHandler struct {
	fundService *services.FundService
}

func NewFundHandler(fundService *services.FundService) *FundHandler {
	return &FundHandler{
		fundService: fundService,
	}
}

// GetAllFunds retrieves all funds for the authenticated user
// GET /funds
func (fh *FundHandler) GetAllFunds(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	funds, err := fh.fundService.GetAllFunds(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Build response with computed values
	fundResponses := make([]models.FundResponse, 0, len(funds))
	for _, fund := range funds {
		totalPaid, err := fh.fundService.CalculateTotalPaid(c.Context(), fund.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		outstanding, err := fh.fundService.CalculateOutstanding(c.Context(), &fund)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		status, err := fh.fundService.GetFundStatus(c.Context(), &fund)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		transactions, err := fh.fundService.GetTransactionsByFundID(c.Context(), fund.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		fundResponses = append(fundResponses, models.FundResponse{
			ID:              fund.ID.Hex(),
			PersonName:      fund.PersonName,
			Type:            fund.Type,
			PrincipalAmount: fund.PrincipalAmount,
			StartDate:       fund.StartDate,
			Notes:           fund.Notes,
			TotalPaid:       totalPaid,
			Outstanding:     outstanding,
			Status:          status,
			Transactions:    transactions,
			CreatedAt:       fund.CreatedAt,
			UpdatedAt:       fund.UpdatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fundResponses)
}

// GetFundByID retrieves a specific fund by ID
// GET /funds/:fundId
func (fh *FundHandler) GetFundByID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	fundID := c.Params("fundId")

	fund, err := fh.fundService.GetFundByID(c.Context(), userID, fundID)
	if err != nil {
		if err.Error() == "fund not found or doesn't belong to user" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Calculate computed values
	totalPaid, err := fh.fundService.CalculateTotalPaid(c.Context(), fund.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	outstanding, err := fh.fundService.CalculateOutstanding(c.Context(), fund)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	status, err := fh.fundService.GetFundStatus(c.Context(), fund)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	transactions, err := fh.fundService.GetTransactionsByFundID(c.Context(), fund.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.FundResponse{
		ID:              fund.ID.Hex(),
		PersonName:      fund.PersonName,
		Type:            fund.Type,
		PrincipalAmount: fund.PrincipalAmount,
		StartDate:       fund.StartDate,
		Notes:           fund.Notes,
		TotalPaid:       totalPaid,
		Outstanding:     outstanding,
		Status:          status,
		Transactions:    transactions,
		CreatedAt:       fund.CreatedAt,
		UpdatedAt:       fund.UpdatedAt,
	})
}

// CreateFund creates a new fund
// POST /funds
func (fh *FundHandler) CreateFund(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req models.FundRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}

	// Validate required fields
	if req.PersonName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "person name is required",
		})
	}

	if req.Type != models.FundTypeBorrowed && req.Type != models.FundTypeGiven {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "type must be BORROWED or GIVEN",
		})
	}

	if req.PrincipalAmount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "principal amount must be greater than 0",
		})
	}

	fund, err := fh.fundService.CreateFund(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Calculate computed values
	totalPaid, err := fh.fundService.CalculateTotalPaid(c.Context(), fund.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	outstanding, err := fh.fundService.CalculateOutstanding(c.Context(), fund)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	status, err := fh.fundService.GetFundStatus(c.Context(), fund)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.FundResponse{
		ID:              fund.ID.Hex(),
		PersonName:      fund.PersonName,
		Type:            fund.Type,
		PrincipalAmount: fund.PrincipalAmount,
		StartDate:       fund.StartDate,
		Notes:           fund.Notes,
		TotalPaid:       totalPaid,
		Outstanding:     outstanding,
		Status:          status,
		Transactions:    []models.Transaction{},
		CreatedAt:       fund.CreatedAt,
		UpdatedAt:       fund.UpdatedAt,
	})
}

// UpdateFund updates an existing fund
// PUT /funds/:fundId
func (fh *FundHandler) UpdateFund(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	fundID := c.Params("fundId")

	var req models.FundRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}

	// Validate required fields
	if req.PersonName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "person name is required",
		})
	}

	if req.Type != models.FundTypeBorrowed && req.Type != models.FundTypeGiven {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "type must be BORROWED or GIVEN",
		})
	}

	if req.PrincipalAmount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "principal amount must be greater than 0",
		})
	}

	fund, err := fh.fundService.UpdateFund(c.Context(), userID, fundID, req)
	if err != nil {
		if err.Error() == "fund not found or doesn't belong to user" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Calculate computed values
	totalPaid, err := fh.fundService.CalculateTotalPaid(c.Context(), fund.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	outstanding, err := fh.fundService.CalculateOutstanding(c.Context(), fund)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	status, err := fh.fundService.GetFundStatus(c.Context(), fund)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	transactions, err := fh.fundService.GetTransactionsByFundID(c.Context(), fund.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.FundResponse{
		ID:              fund.ID.Hex(),
		PersonName:      fund.PersonName,
		Type:            fund.Type,
		PrincipalAmount: fund.PrincipalAmount,
		StartDate:       fund.StartDate,
		Notes:           fund.Notes,
		TotalPaid:       totalPaid,
		Outstanding:     outstanding,
		Status:          status,
		Transactions:    transactions,
		CreatedAt:       fund.CreatedAt,
		UpdatedAt:       fund.UpdatedAt,
	})
}

// DeleteFund deletes a fund and all its transactions
// DELETE /funds/:fundId
func (fh *FundHandler) DeleteFund(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	fundID := c.Params("fundId")

	err := fh.fundService.DeleteFund(c.Context(), userID, fundID)
	if err != nil {
		if err.Error() == "fund not found or doesn't belong to user" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "fund deleted successfully",
	})
}

// AddTransaction adds a new transaction to a fund
// POST /funds/:fundId/transactions
func (fh *FundHandler) AddTransaction(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	fundID := c.Params("fundId")

	var req models.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}

	// Validate transaction amount
	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "transaction amount must be greater than 0",
		})
	}

	_, err := fh.fundService.AddTransaction(c.Context(), userID, fundID, req)
	if err != nil {
		if err.Error() == "fund not found or doesn't belong to user" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return updated fund with computed values
	fund, err := fh.fundService.GetFundByID(c.Context(), userID, fundID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPaid, err := fh.fundService.CalculateTotalPaid(c.Context(), fund.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	outstanding, err := fh.fundService.CalculateOutstanding(c.Context(), fund)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	status, err := fh.fundService.GetFundStatus(c.Context(), fund)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	transactions, err := fh.fundService.GetTransactionsByFundID(c.Context(), fund.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.FundResponse{
		ID:              fund.ID.Hex(),
		PersonName:      fund.PersonName,
		Type:            fund.Type,
		PrincipalAmount: fund.PrincipalAmount,
		StartDate:       fund.StartDate,
		Notes:           fund.Notes,
		TotalPaid:       totalPaid,
		Outstanding:     outstanding,
		Status:          status,
		Transactions:    transactions,
		CreatedAt:       fund.CreatedAt,
		UpdatedAt:       fund.UpdatedAt,
	})
}

// UpdateTransaction updates an existing transaction
// PUT /funds/:fundId/transactions/:transactionId
func (fh *FundHandler) UpdateTransaction(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	fundID := c.Params("fundId")
	transactionID := c.Params("transactionId")

	var req models.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}

	// Validate transaction amount
	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "transaction amount must be greater than 0",
		})
	}

	_, err := fh.fundService.UpdateTransaction(c.Context(), userID, fundID, transactionID, req)
	if err != nil {
		if err.Error() == "fund not found or doesn't belong to user" || err.Error() == "transaction not found or doesn't belong to fund" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return updated fund with computed values
	fund, err := fh.fundService.GetFundByID(c.Context(), userID, fundID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPaid, err := fh.fundService.CalculateTotalPaid(c.Context(), fund.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	outstanding, err := fh.fundService.CalculateOutstanding(c.Context(), fund)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	status, err := fh.fundService.GetFundStatus(c.Context(), fund)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	transactions, err := fh.fundService.GetTransactionsByFundID(c.Context(), fund.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.FundResponse{
		ID:              fund.ID.Hex(),
		PersonName:      fund.PersonName,
		Type:            fund.Type,
		PrincipalAmount: fund.PrincipalAmount,
		StartDate:       fund.StartDate,
		Notes:           fund.Notes,
		TotalPaid:       totalPaid,
		Outstanding:     outstanding,
		Status:          status,
		Transactions:    transactions,
		CreatedAt:       fund.CreatedAt,
		UpdatedAt:       fund.UpdatedAt,
	})
}

// DeleteTransaction deletes a transaction
// DELETE /funds/:fundId/transactions/:transactionId
func (fh *FundHandler) DeleteTransaction(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	fundID := c.Params("fundId")
	transactionID := c.Params("transactionId")

	err := fh.fundService.DeleteTransaction(c.Context(), userID, fundID, transactionID)
	if err != nil {
		if err.Error() == "fund not found or doesn't belong to user" || err.Error() == "transaction not found or doesn't belong to fund" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return updated fund with computed values
	fund, err := fh.fundService.GetFundByID(c.Context(), userID, fundID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPaid, err := fh.fundService.CalculateTotalPaid(c.Context(), fund.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	outstanding, err := fh.fundService.CalculateOutstanding(c.Context(), fund)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	status, err := fh.fundService.GetFundStatus(c.Context(), fund)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	transactions, err := fh.fundService.GetTransactionsByFundID(c.Context(), fund.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.FundResponse{
		ID:              fund.ID.Hex(),
		PersonName:      fund.PersonName,
		Type:            fund.Type,
		PrincipalAmount: fund.PrincipalAmount,
		StartDate:       fund.StartDate,
		Notes:           fund.Notes,
		TotalPaid:       totalPaid,
		Outstanding:     outstanding,
		Status:          status,
		Transactions:    transactions,
		CreatedAt:       fund.CreatedAt,
		UpdatedAt:       fund.UpdatedAt,
	})
}

