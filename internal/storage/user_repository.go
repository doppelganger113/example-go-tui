package storage

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Id    primitive.ObjectID `bson:"_id"`
	Email string
}

type UserRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) (*UserRepository, error) {
	return &UserRepository{db: db, coll: db.Collection("users")}, nil
}

func (r *UserRepository) AddUserIfNotExists(ctx context.Context, user User) error {
	existing, err := r.FindByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failing searching for existing: %v", err)
	}
	if existing != nil {
		return nil
	}

	_, err = r.coll.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.coll.FindOne(ctx, bson.D{{"email", email}}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed searching user by email: %v", err)
	}

	return &user, nil
}
