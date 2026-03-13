package usecase

import (
	"errors"
	"strings"

	"github.com/philipos/prepbase/domain"
)

type questionUsecase struct {
	questionRepo domain.QuestionRepository
}

func NewQuestionUsecase(qRepo domain.QuestionRepository) domain.QuestionUsecase {
	return &questionUsecase{
		questionRepo: qRepo,
	}
}

func (u *questionUsecase) Create(q *domain.Question) error {
	// Validation
	if strings.TrimSpace(q.Title) == "" {
		return errors.New("title cannot be empty")
	}
	if strings.TrimSpace(q.Content) == "" {
		return errors.New("content cannot be empty")
	}
	if strings.TrimSpace(q.AuthorID) == "" {
		return errors.New("author ID is required")
	}

	// Set initial metrics
	q.Upvotes = 0
	q.Downvotes = 0
	q.Views = 0

	return u.questionRepo.Create(q)
}

func (u *questionUsecase) GetByID(id string) (*domain.Question, error) {
	q, err := u.questionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Fire and forget view increment. 
	// We don't return an error if this fails because the user still got their data.
	// In Phase 5, we will move this to a Goroutine! ???
	_ = u.questionRepo.IncrementViews(id)

	return q, nil
}

func (u *questionUsecase) FetchAll(cursor string, limit int) ([]domain.Question, error) {
	if limit <= 0 || limit > 50 {
		limit = 10 // Default fallback
	}
	return u.questionRepo.FetchAll(cursor, limit)
}

func (u *questionUsecase) Update(id string, authorID string, q *domain.Question) error {
	// Verify the question exists
	existingQ, err := u.questionRepo.GetByID(id)
	if err != nil {
		return errors.New("question not found")
	}

	// AUTHORIZATION: Only the author can update it
	if existingQ.AuthorID != authorID {
		return errors.New("unauthorized: you can only update your own questions")
	}

	// Validation
	if strings.TrimSpace(q.Title) == "" {
		return errors.New("title cannot be empty")
	}

	q.ID = id
	return u.questionRepo.Update(q)
}

func (u *questionUsecase) Delete(id string, authorID string, userRole string) error {
	// Verify the question exists
	existingQ, err := u.questionRepo.GetByID(id)
	if err != nil {
		return errors.New("question not found")
	}

	// 2. AUTHORIZATION: Author OR Admin can delete
	if existingQ.AuthorID != authorID && userRole != "admin" {
		return errors.New("unauthorized: you do not have permission to delete this question")
	}

	return u.questionRepo.Delete(id)
}

// ----------------------------------------------------
// PLACEHOLDERS FOR PHASE 4 (Engagement) ???
// ----------------------------------------------------

func (u *questionUsecase) Upvote(userID string, questionID string) error {
	return errors.New("upvote logic not yet implemented")
}

func (u *questionUsecase) Downvote(userID string, questionID string) error {
	return errors.New("downvote logic not yet implemented")
}