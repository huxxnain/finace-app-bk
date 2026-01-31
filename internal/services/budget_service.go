package services

import (
	"context"
	"fmt"
	"time"

	"github.com/huxxnainali/finance-app/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BudgetService struct {
	collection *mongo.Collection
}

func NewBudgetService(db *mongo.Database) *BudgetService {
	collection := db.Collection("monthly_budgets")

	// Create unique index on userId, year, month
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "userId", Value: 1},
			{Key: "year", Value: 1},
			{Key: "month", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	collection.Indexes().CreateOne(context.Background(), indexModel)

	return &BudgetService{collection: collection}
}

// GetOrCreateBudget retrieves a budget or creates one if it doesn't exist
func (bs *BudgetService) GetOrCreateBudget(ctx context.Context, userID string, year, month int) (*models.MonthlyBudget, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	budget := &models.MonthlyBudget{}
	err = bs.collection.FindOne(ctx, bson.M{
		"userId": objID,
		"year":   year,
		"month":  month,
	}).Decode(budget)

	// If budget doesn't exist, create it
	if err == mongo.ErrNoDocuments {
		budget = &models.MonthlyBudget{
			ID:         primitive.NewObjectID(),
			UserID:     objID,
			Year:       year,
			Month:      month,
			BaseIncome: nil,
			Expenses:   []models.Expense{},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		_, err := bs.collection.InsertOne(ctx, budget)
		if err != nil {
			return nil, err
		}

		return budget, nil
	}

	if err != nil {
		return nil, err
	}

	return budget, nil
}

// SetBaseIncome sets or updates the base income for a month
func (bs *BudgetService) SetBaseIncome(ctx context.Context, userID string, year, month int, amount float64) (*models.MonthlyBudget, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Update base income
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := &models.MonthlyBudget{}
	err = bs.collection.FindOneAndUpdate(ctx,
		bson.M{
			"userId": objID,
			"year":   year,
			"month":  month,
		},
		bson.M{
			"$set": bson.M{
				"baseIncome": amount,
				"updatedAt":  time.Now(),
			},
		},
		opts,
	).Decode(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// AddExpense adds an expense to a budget
func (bs *BudgetService) AddExpense(ctx context.Context, userID string, year, month int, expense models.Expense) (*models.MonthlyBudget, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Ensure expense has an ID
	if expense.ID == primitive.NilObjectID {
		expense.ID = primitive.NewObjectID()
	}
	if expense.CreatedAt.IsZero() {
		expense.CreatedAt = time.Now()
	}

	// Get or create budget
	_, err = bs.GetOrCreateBudget(ctx, userID, year, month)
	if err != nil {
		return nil, err
	}

	// Add expense
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := &models.MonthlyBudget{}
	err = bs.collection.FindOneAndUpdate(ctx,
		bson.M{
			"userId": objID,
			"year":   year,
			"month":  month,
		},
		bson.M{
			"$push": bson.M{
				"expenses": expense,
			},
			"$set": bson.M{
				"updatedAt": time.Now(),
			},
		},
		opts,
	).Decode(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateExpense updates an existing expense
func (bs *BudgetService) UpdateExpense(ctx context.Context, userID, expenseID string, updatedExpense models.Expense) (*models.MonthlyBudget, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	expenseObjID, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		return nil, fmt.Errorf("invalid expense ID")
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := &models.MonthlyBudget{}
	err = bs.collection.FindOneAndUpdate(ctx,
		bson.M{
			"userId":       objID,
			"expenses._id": expenseObjID,
		},
		bson.M{
			"$set": bson.M{
				"expenses.$.title":  updatedExpense.Title,
				"expenses.$.amount": updatedExpense.Amount,
				"updatedAt":         time.Now(),
			},
		},
		opts,
	).Decode(result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("expense not found or doesn't belong to user")
		}
		return nil, err
	}

	return result, nil
}

// DeleteExpense deletes an expense from a budget
func (bs *BudgetService) DeleteExpense(ctx context.Context, userID, expenseID string) (*models.MonthlyBudget, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	expenseObjID, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		return nil, fmt.Errorf("invalid expense ID")
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := &models.MonthlyBudget{}
	err = bs.collection.FindOneAndUpdate(ctx,
		bson.M{
			"userId":       objID,
			"expenses._id": expenseObjID,
		},
		bson.M{
			"$pull": bson.M{
				"expenses": bson.M{"_id": expenseObjID},
			},
			"$set": bson.M{
				"updatedAt": time.Now(),
			},
		},
		opts,
	).Decode(result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("expense not found or doesn't belong to user")
		}
		return nil, err
	}

	return result, nil
}

// CalculateRemaining calculates the remaining balance
func CalculateRemaining(baseIncome *float64, expenses []models.Expense) *float64 {
	if baseIncome == nil {
		return nil
	}

	total := *baseIncome
	for _, expense := range expenses {
		total -= expense.Amount
	}

	return &total
}
