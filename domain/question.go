package domain

import "time"

type Question struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`    
	Experience string    `json:"experience"` 
	Difficulty string    `json:"difficulty"` 
	Company    string    `json:"company"` 
	Tags       []string  `json:"tags"`
	Upvotes    int       `json:"upvotes"`
	Downvotes  int       `json:"downvotes"`
	Views      int       `json:"views"`
	AuthorID   string    `json:"author_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Vote struct {
	UserID     string `json:"user_id"`
	QuestionID string `json:"question_id"`
	Value      int    `json:"value"` 
}

type VoteRepository interface {
	GetVote(userID string, questionID string) (*Vote, error)
	AddVote(vote *Vote) error
	UpdateVote(vote *Vote) error
	DeleteVote(userID string, questionID string) error
}


type QuestionRepository interface {
	Create(question *Question) error
	GetByID(id string) (*Question, error)
	FetchAll(cursor string, limit int) ([]Question, error)
	Update(question *Question) error
	Delete(id string) error
	UpdateVoteCount(questionID string, upvoteChange int, downvoteChange int) error
	IncrementViews(id string) error
}

type QuestionUsecase interface {
	Create(question *Question) error
	GetByID(id string) (*Question, error)
	FetchAll(cursor string, limit int) ([]Question, error)
	Update(id string, authorID string, question *Question) error
	Delete(id string, authorID string, userRole string) error

	Upvote(userID string, questionID string) error
	Downvote(userID string, questionID string) error
	GenerateAIHint(id string) (string, error)
}