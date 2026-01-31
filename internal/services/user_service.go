package services

import (
	"context"
	"fmt"
	"time"

	"github.com/huxxnainali/finance-app/internal/models"
	"github.com/huxxnainali/finance-app/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService struct {
	collection *mongo.Collection
}

func NewUserService(db *mongo.Database) *UserService {
	collection := db.Collection("users_expense")

	// Create unique index on email
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	collection.Indexes().CreateOne(context.Background(), indexModel)

	return &UserService{collection: collection}
}

// SignUp creates a new user with hashed password
func (us *UserService) SignUp(ctx context.Context, email, password string) (*models.User, error) {
	// Check if user already exists
	existingUser := &models.User{}
	err := us.collection.FindOne(ctx, bson.M{"email": email}).Decode(existingUser)
	if err == nil {
		return nil, fmt.Errorf("user with this email already exists")
	}
	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &models.User{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	_, err = us.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns the user if credentials are valid
func (us *UserService) Login(ctx context.Context, email, password string) (*models.User, error) {
	user := &models.User{}
	err := us.collection.FindOne(ctx, bson.M{"email": email}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, err
	}

	// Compare passwords
	err = utils.ComparePasswords(user.Password, password)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (us *UserService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	user := &models.User{}
	err = us.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}
