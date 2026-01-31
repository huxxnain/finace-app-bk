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
	Year   int     `json:"year"`
	Month  int     `json:"month"`
}

// ExpenseRequest is the request format for expense endpoints
type ExpenseRequest struct {
	Title  string  `json:"title"`
	Amount float64 `json:"amount"`
	Year   int     `json:"year"`
	Month  int     `json:"month"`
}

// JWTClaims represents JWT claims
type JWTClaims struct {
	UserID string `json:"userId"`
}

// FundType represents the type of fund
type FundType string

const (
	FundTypeBorrowed FundType = "BORROWED"
	FundTypeGiven    FundType = "GIVEN"
)

// Fund represents a borrowing or lending agreement
type Fund struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	PersonName     string             `bson:"personName" json:"personName"`
	Type           FundType           `bson:"type" json:"type"`
	PrincipalAmount float64           `bson:"principalAmount" json:"principalAmount"`
	StartDate      time.Time          `bson:"startDate" json:"startDate"`
	Notes          string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// Transaction represents a partial payment for a fund
type Transaction struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FundID    primitive.ObjectID `bson:"fundId" json:"fundId"`
	Amount    float64            `bson:"amount" json:"amount"`
	Date      time.Time          `bson:"date" json:"date"`
	Note      string             `bson:"note,omitempty" json:"note,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

// FundRequest is the request format for fund endpoints
type FundRequest struct {
	PersonName      string    `json:"personName"`
	Type            FundType  `json:"type"`
	PrincipalAmount float64   `json:"principalAmount"`
	StartDate       time.Time `json:"startDate"`
	Notes           string    `json:"notes,omitempty"`
}

// TransactionRequest is the request format for transaction endpoints
type TransactionRequest struct {
	Amount float64   `json:"amount"`
	Date   time.Time `json:"date"`
	Note   string    `json:"note,omitempty"`
}

// FundResponse is the response format for fund endpoints
type FundResponse struct {
	ID              string       `json:"id"`
	PersonName      string       `json:"personName"`
	Type            FundType     `json:"type"`
	PrincipalAmount float64      `json:"principalAmount"`
	StartDate       time.Time    `json:"startDate"`
	Notes           string       `json:"notes,omitempty"`
	TotalPaid       float64      `json:"totalPaid"`
	Outstanding     float64      `json:"outstanding"`
	Status          string       `json:"status"`
	Transactions    []Transaction `json:"transactions"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
}