package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

// Expense represents a single expense
type Expense struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Amount    float64            `bson:"amount" json:"amount"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

// MonthlyBudget represents a user's budget for a specific month
type MonthlyBudget struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"userId" json:"userId"`
	Year       int                `bson:"year" json:"year"`
	Month      int                `bson:"month" json:"month"`
	BaseIncome *float64           `bson:"baseIncome" json:"baseIncome"`
	Expenses   []Expense          `bson:"expenses" json:"expenses"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// BudgetResponse is the response format for budget endpoints
type BudgetResponse struct {
	Year       int       `json:"year"`
	Month      int       `json:"month"`
	BaseIncome *float64  `json:"baseIncome"`
	Expenses   []Expense `json:"expenses"`
	Remaining  *float64  `json:"remaining"`
}

// AuthRequest is the request format for auth endpoints
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse is the response format for auth endpoints
type AuthResponse struct {
	Token string `json:"token"`
}

// BaseIncomeRequest is the request format for setting base income
type BaseIncomeRequest struct {
	Amount float64 `json:"amount"`
}

// ExpenseRequest is the request format for expense endpoints
type ExpenseRequest struct {
	Title  string  `json:"title"`
	Amount float64 `json:"amount"`
}

// JWTClaims represents JWT claims
type JWTClaims struct {
	UserID string `json:"userId"`
}
