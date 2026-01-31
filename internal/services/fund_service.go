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

type FundService struct {
	fundCollection        *mongo.Collection
	transactionCollection *mongo.Collection
}

func NewFundService(db *mongo.Database) *FundService {
	fundCollection := db.Collection("funds")
	transactionCollection := db.Collection("transactions")

	// Create index on userId for funds
	fundIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "userId", Value: 1},
		},
	}
	fundCollection.Indexes().CreateOne(context.Background(), fundIndexModel)

	// Create index on fundId for transactions
	transactionIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "fundId", Value: 1},
		},
	}
	transactionCollection.Indexes().CreateOne(context.Background(), transactionIndexModel)

	return &FundService{
		fundCollection:        fundCollection,
		transactionCollection: transactionCollection,
	}
}

// CalculateTotalPaid calculates the total amount paid from all transactions
func (fs *FundService) CalculateTotalPaid(ctx context.Context, fundID primitive.ObjectID) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"fundId": fundID}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}}},
	}

	cursor, err := fs.transactionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Total float64 `bson:"total"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.Total, nil
	}

	return 0, nil
}

// CalculateOutstanding calculates outstanding amount (principalAmount - totalPaid)
func (fs *FundService) CalculateOutstanding(ctx context.Context, fund *models.Fund) (float64, error) {
	totalPaid, err := fs.CalculateTotalPaid(ctx, fund.ID)
	if err != nil {
		return 0, err
	}

	outstanding := fund.PrincipalAmount - totalPaid
	if outstanding < 0 {
		outstanding = 0
	}

	return outstanding, nil
}

// GetFundStatus returns "OPEN" if outstanding > 0, "PAID" if outstanding === 0
func (fs *FundService) GetFundStatus(ctx context.Context, fund *models.Fund) (string, error) {
	outstanding, err := fs.CalculateOutstanding(ctx, fund)
	if err != nil {
		return "", err
	}

	if outstanding > 0 {
		return "OPEN", nil
	}
	return "PAID", nil
}

// GetAllFunds retrieves all funds for a user
func (fs *FundService) GetAllFunds(ctx context.Context, userID string) ([]models.Fund, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	cursor, err := fs.fundCollection.Find(ctx, bson.M{"userId": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var funds []models.Fund
	if err := cursor.All(ctx, &funds); err != nil {
		return nil, err
	}

	return funds, nil
}

// GetFundByID retrieves a fund by ID
func (fs *FundService) GetFundByID(ctx context.Context, userID, fundID string) (*models.Fund, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	fundObjID, err := primitive.ObjectIDFromHex(fundID)
	if err != nil {
		return nil, fmt.Errorf("invalid fund ID")
	}

	fund := &models.Fund{}
	err = fs.fundCollection.FindOne(ctx, bson.M{
		"_id":    fundObjID,
		"userId": userObjID,
	}).Decode(fund)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("fund not found or doesn't belong to user")
		}
		return nil, err
	}

	return fund, nil
}

// CreateFund creates a new fund
func (fs *FundService) CreateFund(ctx context.Context, userID string, req models.FundRequest) (*models.Fund, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Validate fund type
	if req.Type != models.FundTypeBorrowed && req.Type != models.FundTypeGiven {
		return nil, fmt.Errorf("invalid fund type, must be BORROWED or GIVEN")
	}

	// Validate principal amount
	if req.PrincipalAmount <= 0 {
		return nil, fmt.Errorf("principal amount must be greater than 0")
	}

	now := time.Now()
	fund := &models.Fund{
		ID:              primitive.NewObjectID(),
		UserID:          objID,
		PersonName:      req.PersonName,
		Type:            req.Type,
		PrincipalAmount: req.PrincipalAmount,
		StartDate:       req.StartDate,
		Notes:           req.Notes,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	_, err = fs.fundCollection.InsertOne(ctx, fund)
	if err != nil {
		return nil, err
	}

	return fund, nil
}

// UpdateFund updates an existing fund
func (fs *FundService) UpdateFund(ctx context.Context, userID, fundID string, req models.FundRequest) (*models.Fund, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	fundObjID, err := primitive.ObjectIDFromHex(fundID)
	if err != nil {
		return nil, fmt.Errorf("invalid fund ID")
	}

	// Verify fund exists and belongs to user
	_, err = fs.GetFundByID(ctx, userID, fundID)
	if err != nil {
		return nil, err
	}

	// Validate fund type
	if req.Type != models.FundTypeBorrowed && req.Type != models.FundTypeGiven {
		return nil, fmt.Errorf("invalid fund type, must be BORROWED or GIVEN")
	}

	// Calculate current total paid
	totalPaid, err := fs.CalculateTotalPaid(ctx, fundObjID)
	if err != nil {
		return nil, err
	}

	// Validate: principalAmount MUST be >= totalPaid
	if req.PrincipalAmount < totalPaid {
		return nil, fmt.Errorf("principal amount cannot be less than total paid (%.2f)", totalPaid)
	}

	// Validate principal amount
	if req.PrincipalAmount <= 0 {
		return nil, fmt.Errorf("principal amount must be greater than 0")
	}

	// Update fund
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := &models.Fund{}
	err = fs.fundCollection.FindOneAndUpdate(ctx,
		bson.M{
			"_id":    fundObjID,
			"userId": userObjID,
		},
		bson.M{
			"$set": bson.M{
				"personName":      req.PersonName,
				"type":            req.Type,
				"principalAmount": req.PrincipalAmount,
				"startDate":       req.StartDate,
				"notes":           req.Notes,
				"updatedAt":       time.Now(),
			},
		},
		opts,
	).Decode(result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("fund not found or doesn't belong to user")
		}
		return nil, err
	}

	return result, nil
}

// DeleteFund deletes a fund and all its transactions
func (fs *FundService) DeleteFund(ctx context.Context, userID, fundID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	fundObjID, err := primitive.ObjectIDFromHex(fundID)
	if err != nil {
		return fmt.Errorf("invalid fund ID")
	}

	// Verify fund belongs to user
	_, err = fs.GetFundByID(ctx, userID, fundID)
	if err != nil {
		return err
	}

	// Delete all transactions for this fund
	_, err = fs.transactionCollection.DeleteMany(ctx, bson.M{"fundId": fundObjID})
	if err != nil {
		return err
	}

	// Delete fund
	_, err = fs.fundCollection.DeleteOne(ctx, bson.M{
		"_id":    fundObjID,
		"userId": userObjID,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetTransactionsByFundID retrieves all transactions for a fund
func (fs *FundService) GetTransactionsByFundID(ctx context.Context, fundID primitive.ObjectID) ([]models.Transaction, error) {
	cursor, err := fs.transactionCollection.Find(ctx, bson.M{"fundId": fundID}, options.Find().SetSort(bson.M{"date": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []models.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

// AddTransaction adds a new transaction to a fund
func (fs *FundService) AddTransaction(ctx context.Context, userID, fundID string, req models.TransactionRequest) (*models.Transaction, error) {
	fundObjID, err := primitive.ObjectIDFromHex(fundID)
	if err != nil {
		return nil, fmt.Errorf("invalid fund ID")
	}

	// Verify fund belongs to user
	fund, err := fs.GetFundByID(ctx, userID, fundID)
	if err != nil {
		return nil, err
	}

	// Validate transaction amount
	if req.Amount <= 0 {
		return nil, fmt.Errorf("transaction amount must be greater than 0")
	}

	// Calculate current total paid
	totalPaid, err := fs.CalculateTotalPaid(ctx, fundObjID)
	if err != nil {
		return nil, err
	}

	// Validate: sum(transactions.amount) <= principalAmount
	if totalPaid+req.Amount > fund.PrincipalAmount {
		return nil, fmt.Errorf("transaction amount would exceed principal amount. Maximum allowed: %.2f", fund.PrincipalAmount-totalPaid)
	}

	// Create transaction
	transaction := &models.Transaction{
		ID:        primitive.NewObjectID(),
		FundID:    fundObjID,
		Amount:    req.Amount,
		Date:      req.Date,
		Note:      req.Note,
		CreatedAt: time.Now(),
	}

	_, err = fs.transactionCollection.InsertOne(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// UpdateTransaction updates an existing transaction
func (fs *FundService) UpdateTransaction(ctx context.Context, userID, fundID, transactionID string, req models.TransactionRequest) (*models.Transaction, error) {
	fundObjID, err := primitive.ObjectIDFromHex(fundID)
	if err != nil {
		return nil, fmt.Errorf("invalid fund ID")
	}

	transactionObjID, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction ID")
	}

	// Verify fund belongs to user
	fund, err := fs.GetFundByID(ctx, userID, fundID)
	if err != nil {
		return nil, err
	}

	// Verify transaction belongs to fund
	existingTransaction := &models.Transaction{}
	err = fs.transactionCollection.FindOne(ctx, bson.M{
		"_id":    transactionObjID,
		"fundId": fundObjID,
	}).Decode(existingTransaction)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("transaction not found or doesn't belong to fund")
		}
		return nil, err
	}

	// Validate transaction amount
	if req.Amount <= 0 {
		return nil, fmt.Errorf("transaction amount must be greater than 0")
	}

	// Calculate current total paid (excluding the transaction being updated)
	totalPaid, err := fs.CalculateTotalPaid(ctx, fundObjID)
	if err != nil {
		return nil, err
	}
	totalPaidWithoutThis := totalPaid - existingTransaction.Amount

	// Validate: sum(transactions.amount) <= principalAmount
	if totalPaidWithoutThis+req.Amount > fund.PrincipalAmount {
		return nil, fmt.Errorf("transaction amount would exceed principal amount. Maximum allowed: %.2f", fund.PrincipalAmount-totalPaidWithoutThis)
	}

	// Update transaction
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := &models.Transaction{}
	err = fs.transactionCollection.FindOneAndUpdate(ctx,
		bson.M{
			"_id":    transactionObjID,
			"fundId": fundObjID,
		},
		bson.M{
			"$set": bson.M{
				"amount": req.Amount,
				"date":   req.Date,
				"note":   req.Note,
			},
		},
		opts,
	).Decode(result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("transaction not found or doesn't belong to fund")
		}
		return nil, err
	}

	return result, nil
}

// DeleteTransaction deletes a transaction
func (fs *FundService) DeleteTransaction(ctx context.Context, userID, fundID, transactionID string) error {
	fundObjID, err := primitive.ObjectIDFromHex(fundID)
	if err != nil {
		return fmt.Errorf("invalid fund ID")
	}

	transactionObjID, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		return fmt.Errorf("invalid transaction ID")
	}

	// Verify fund belongs to user
	_, err = fs.GetFundByID(ctx, userID, fundID)
	if err != nil {
		return err
	}

	// Verify transaction belongs to fund
	var transaction models.Transaction
	err = fs.transactionCollection.FindOne(ctx, bson.M{
		"_id":    transactionObjID,
		"fundId": fundObjID,
	}).Decode(&transaction)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("transaction not found or doesn't belong to fund")
		}
		return err
	}

	// Delete transaction
	_, err = fs.transactionCollection.DeleteOne(ctx, bson.M{
		"_id":    transactionObjID,
		"fundId": fundObjID,
	})
	if err != nil {
		return err
	}

	return nil
}

