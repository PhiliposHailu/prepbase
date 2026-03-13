package repository

import (
	"context"
	"errors"
	"time"

	"github.com/philipos/prepbase/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// 1. THE DATABASE MODEL & MAPPING

// Private struct: Only MongoDB sees this
type userModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Bio       string             `bson:"bio"`
	Role      string             `bson:"role"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"` // Soft Delete
}

// Domain -> DB
func fromUserDomain(u *domain.User) userModel {
	objID, _ := primitive.ObjectIDFromHex(u.ID)
	return userModel{
		ID:        objID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Bio:       u.Bio,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}
}

// DB -> Domain
func toUserDomain(m userModel) domain.User {
	return domain.User{
		ID:        m.ID.Hex(),
		Username:  m.Username,
		Email:     m.Email,
		Password:  m.Password,
		Bio:       m.Bio,
		Role:      m.Role,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}

// 2. THE REPOSITORY IMPLEMENTATION

type userRepository struct {
	database   *mongo.Database
	collection string
}

func NewUserRepository(db *mongo.Database, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

func (r *userRepository) Create(user *domain.User) error {
	collection := r.database.Collection(r.collection)

	// Set timestamps before saving
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	dbUser := fromUserDomain(user)

	res, err := collection.InsertOne(context.TODO(), dbUser)
	if err != nil {
		return err
	}

	user.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	collection := r.database.Collection(r.collection)
	var dbUser userModel

	// Find the user by Email
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&dbUser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	domainUser := toUserDomain(dbUser)
	return &domainUser, nil
}

func (r *userRepository) GetByID(id string) (*domain.User, error) {
	collection := r.database.Collection(r.collection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	var dbUser userModel
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&dbUser)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	domainUser := toUserDomain(dbUser)
	return &domainUser, nil
}

func (r *userRepository) Update(user *domain.User) error {
	collection := r.database.Collection(r.collection)
	objID, _ := primitive.ObjectIDFromHex(user.ID)

	updateFileds := bson.M{}

	if user.Username != "" {
		updateFileds["username"] = user.Username
	}

	if user.Bio != "" {
		updateFileds["bio"] = user.Bio
	}

	if user.Bio != "" {
		updateFileds["role"] = user.Role
	}

	if len(updateFileds) > 0 {
		updateFileds["updated_at"] = time.Now()
	} else {
		return nil
	}

	update := bson.M{
		"$set": updateFileds,
	}

	user.UpdatedAt = time.Now()
	
	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	return err
}

func (r *userRepository) Delete(id string) error {
	collection := r.database.Collection(r.collection)
	objID, _ := primitive.ObjectIDFromHex(id)

	// SOFT DELETE: We do not use DeleteOne! We use UpdateOne to set DeletedAt.
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"deleted_at": &now, 
			"updated_at": now,
		},
	}

	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	return err
}