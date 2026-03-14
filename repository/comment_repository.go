package repository

import (
	"context"
	"time"

	"github.com/philipos/prepbase/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type commentModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	QuestionID string             `bson:"question_id"`
	AuthorID   string             `bson:"author_id"`
	Content    string             `bson:"content"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

func fromCommentDomain(c *domain.Comment) commentModel {
	objID, _ := primitive.ObjectIDFromHex(c.ID)
	return commentModel{
		ID:         objID,
		QuestionID: c.QuestionID,
		AuthorID:   c.AuthorID,
		Content:    c.Content,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

func toCommentDomain(m commentModel) domain.Comment {
	return domain.Comment{
		ID:         m.ID.Hex(),
		QuestionID: m.QuestionID,
		AuthorID:   m.AuthorID,
		Content:    m.Content,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

type commentRepository struct {
	db   *mongo.Database
	coll string
}

func NewCommentRepository(db *mongo.Database, collection string) domain.CommentRepository {
	return &commentRepository{db: db, coll: collection}
}

func (r *commentRepository) Create(c *domain.Comment) error {
	collection := r.db.Collection(r.coll)
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	dbModel := fromCommentDomain(c)
	res, err := collection.InsertOne(context.TODO(), dbModel)
	if err != nil {
		return err
	}
	c.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *commentRepository) GetByQuestionID(questionID string) ([]domain.Comment, error) {
	collection := r.db.Collection(r.coll)
	var dbModels []commentModel

	// Sort comments by newest first
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	
	cur, err := collection.Find(context.TODO(), bson.M{"question_id": questionID}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	if err = cur.All(context.TODO(), &dbModels); err != nil {
		return nil, err
	}

	domainComments := make([]domain.Comment, len(dbModels))
	for i, m := range dbModels { // inefficent right ???
		domainComments[i] = toCommentDomain(m)
	}

	return domainComments, nil
}

func (r *commentRepository) GetByID(id string) (*domain.Comment, error) {
	collection := r.db.Collection(r.coll)
	objID, _ := primitive.ObjectIDFromHex(id)

	var m commentModel
	err := collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&m)
	if err != nil {
		return nil, err
	}

	domainC := toCommentDomain(m)
	return &domainC, nil
}

func (r *commentRepository) Update(c *domain.Comment) error {
	collection := r.db.Collection(r.coll)
	objID, _ := primitive.ObjectIDFromHex(c.ID)

	c.UpdatedAt = time.Now()
	update := bson.M{"$set": bson.M{"content": c.Content, "updated_at": c.UpdatedAt}}

	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	return err
}

func (r *commentRepository) Delete(id string) error {
	collection := r.db.Collection(r.coll)
	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	return err
}