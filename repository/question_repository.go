package repository

import (
	"context"
	"errors"
	"time"

	"github.com/philipos/prepbase/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 1. DATABASE MODEL & MAPPING

type questionModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Title      string             `bson:"title"`
	Content    string             `bson:"content"`
	Experience string             `bson:"experience"`
	Difficulty string             `bson:"difficulty"`
	Company    string             `bson:"company"`
	Tags       []string           `bson:"tags"`
	Upvotes    int                `bson:"upvotes"`
	Downvotes  int                `bson:"downvotes"`
	Views      int                `bson:"views"`
	AuthorID   string             `bson:"author_id"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

func fromQuestionDomain(q *domain.Question) questionModel {
	objID, _ := primitive.ObjectIDFromHex(q.ID)
	return questionModel{
		ID:         objID,
		Title:      q.Title,
		Content:    q.Content,
		Experience: q.Experience,
		Difficulty: q.Difficulty,
		Company:    q.Company,
		Tags:       q.Tags,
		Upvotes:    q.Upvotes,
		Downvotes:  q.Downvotes,
		Views:      q.Views,
		AuthorID:   q.AuthorID,
		CreatedAt:  q.CreatedAt,
		UpdatedAt:  q.UpdatedAt,
	}
}

func toQuestionDomain(m questionModel) domain.Question {
	return domain.Question{
		ID:         m.ID.Hex(),
		Title:      m.Title,
		Content:    m.Content,
		Experience: m.Experience,
		Difficulty: m.Difficulty,
		Company:    m.Company,
		Tags:       m.Tags,
		Upvotes:    m.Upvotes,
		Downvotes:  m.Downvotes,
		Views:      m.Views,
		AuthorID:   m.AuthorID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

// 2. REPOSITORY IMPLEMENTATION

type questionRepository struct {
	database   *mongo.Database
	collection string
}

func NewQuestionRepository(db *mongo.Database, collection string) domain.QuestionRepository {
	return &questionRepository{
		database:   db,
		collection: collection,
	}
}

func (r *questionRepository) Create(q *domain.Question) error {
	collection := r.database.Collection(r.collection)

	q.CreatedAt = time.Now()
	q.UpdatedAt = time.Now()

	dbModel := fromQuestionDomain(q)

	res, err := collection.InsertOne(context.TODO(), dbModel)
	if err != nil {
		return err
	}

	q.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *questionRepository) GetByID(id string) (*domain.Question, error) {
	collection := r.database.Collection(r.collection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	var m questionModel
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&m)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("question not found")
	} else if err != nil {
		return nil, err
	}

	domainQ := toQuestionDomain(m)
	return &domainQ, nil
}

func (r *questionRepository) FetchAll(cursor string, limit int) ([]domain.Question, error) {
	collection := r.database.Collection(r.collection)
	var dbModels []questionModel

	// Sort by newest first
	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))

	filter := bson.M{}
	if cursor != "" {
		// Try to parse the cursor as a timestamp (RFC3339 format)
		cursorTime, err := time.Parse(time.RFC3339, cursor)
		if err == nil {
			// Find documents created BEFORE the cursor time (since we sort descending)
			filter["created_at"] = bson.M{"$lt": cursorTime}
		}
	}

	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	if err = cur.All(context.TODO(), &dbModels); err != nil {
		return nil, err
	}

	domainQuestions := make([]domain.Question, len(dbModels))
	for i, m := range dbModels {
		domainQuestions[i] = toQuestionDomain(m)
	}

	return domainQuestions, nil
}

func (r *questionRepository) Update(q *domain.Question) error {
	collection := r.database.Collection(r.collection)
	objID, _ := primitive.ObjectIDFromHex(q.ID)

	// update := bson.M{
	// 	"$set": bson.M{
	// 		"title":      q.Title,
	// 		"content":    q.Content,
	// 		"experience": q.Experience,
	// 		"difficulty": q.Difficulty,
	// 		"company":    q.Company,
	// 		"tags":       q.Tags,
	// 		"updated_at": q.UpdatedAt,
	// 	},
	// }

	updateFileds := bson.M{}

	if q.Title != "" {
		updateFileds["title"] = q.Title
	}
	if q.Content != "" {
		updateFileds["content"] = q.Content
	}
	if q.Experience != "" {
		updateFileds["experience"] = q.Experience
	}
	if q.Difficulty != "" {
		updateFileds["difficulty"] = q.Difficulty
	}
	if q.Company != "" {
		updateFileds["company"] = q.Company
	}
	if len(q.Tags) > 0 {
		updateFileds["tags"] = q.Tags
	}

	if len(updateFileds) > 0 {
		updateFileds["updated_at"] = time.Now()
		q.UpdatedAt = time.Now()

	} else {
		return nil
	}

	update := bson.M{
		"$set": updateFileds,
	}

	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	return err
}

func (r *questionRepository) Delete(id string) error {
	collection := r.database.Collection(r.collection)
	objID, _ := primitive.ObjectIDFromHex(id)

	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	return err
}

func (r *questionRepository) UpdateVoteCount(questionID string, upvoteChange int, downvoteChange int) error {
	collection := r.database.Collection(r.collection)
	objID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		return errors.New("invalid question ID")
	}

	// We use $inc to add or subtract from the current total atomically , oue usecase logic hadles whether it increases or not.
	update := bson.M{
		"$inc": bson.M{
			"upvotes":   upvoteChange,
			"downvotes": downvoteChange,
		},
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	return err
}

func (r *questionRepository) IncrementViews(id string) error {
	collection := r.database.Collection(r.collection)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// $inc is a MongoDB operator that atomically adds a number to a field
	update := bson.M{"$inc": bson.M{"views": 1}}
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	return err
}
