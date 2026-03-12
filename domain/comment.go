package domain

import "time"

type Comment struct {
	ID         string    `json:"id"`
	QuestionID string    `json:"question_id"`
	AuthorID   string    `json:"author_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CommentRepository interface {
	Create(comment *Comment) error
	GetByQuestionID(questionID string) ([]Comment, error)
	GetByID(id string) (*Comment, error)
	Update(comment *Comment) error
	Delete(id string) error
}

type CommentUsecase interface {
	Create(comment *Comment) error
	GetByQuestionID(questionID string) ([]Comment, error)
	Update(id string, authorID string, content string) error
	Delete(id string, authorID string, userRole string) error
}