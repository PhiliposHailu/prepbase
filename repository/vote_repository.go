package repository

import (
	"context"
	"errors"

	"github.com/philipos/prepbase/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// 1. DATABASE MODEL & MAPPING
type voteModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string             `bson:"user_id"`
	QuestionID string             `bson:"question_id"`
	Value      int                `bson:"value"`
}

func fromVoteDomain(v *domain.Vote) voteModel {
	return voteModel{
		UserID:     v.UserID,
		QuestionID: v.QuestionID,
		Value:      v.Value,
	}
}

func toVoteDomain(m voteModel) domain.Vote {
	return domain.Vote{
		UserID:     m.UserID,
		QuestionID: m.QuestionID,
		Value:      m.Value,
	}
}

// 2. REPOSITORY IMPLEMENTATION
type voteRepository struct {
	db   *mongo.Database
	coll string
}

func NewVoteRepository(db *mongo.Database, collection string) domain.VoteRepository {
	return &voteRepository{db: db, coll: collection}
}

func (r *voteRepository) GetVote(userID string, questionID string) (*domain.Vote, error) {
	collection := r.db.Collection(r.coll)
	var m voteModel

	filter := bson.M{
		"user_id":     userID,
		"question_id": questionID,
	}

	err := collection.FindOne(context.TODO(), filter).Decode(&m)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("vote not found")
	} else if err != nil {
		return nil, err
	}

	domainVote := toVoteDomain(m)
	return &domainVote, nil
}

func (r *voteRepository) AddVote(vote *domain.Vote) error {
	collection := r.db.Collection(r.coll)
	
	dbModel := fromVoteDomain(vote)
	_, err := collection.InsertOne(context.TODO(), dbModel)
	return err
}

func (r *voteRepository) UpdateVote(vote *domain.Vote) error {
	collection := r.db.Collection(r.coll)
	filter := bson.M{
		"user_id":     vote.UserID,
		"question_id": vote.QuestionID,
	}
	update := bson.M{"$set": bson.M{"value": vote.Value}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *voteRepository) DeleteVote(userID string, questionID string) error {
	collection := r.db.Collection(r.coll)
	filter := bson.M{
		"user_id":     userID,
		"question_id": questionID,
	}
	_, err := collection.DeleteOne(context.TODO(), filter)
	return err
}