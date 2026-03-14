package usecase

import (
	"errors"
	"strings"

	"github.com/philipos/prepbase/domain"
)

type questionUsecase struct {
	questionRepo domain.QuestionRepository
	voteRepo     domain.VoteRepository
	cache        domain.CacheService
}

type Vote struct {
	UserID     string `json:"user_id" bson:"user_id"`
	QuestionID string `json:"question_id" bson:"question_id"`
	Value      int    `json:"value" bson:"value"`
}

func NewQuestionUsecase(qRepo domain.QuestionRepository, vRepo domain.VoteRepository, cache domain.CacheService) domain.QuestionUsecase {
	return &questionUsecase{
		questionRepo: qRepo,
		voteRepo:     vRepo,
		cache:        cache,
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
	cacheKey := "question_" + id

	// 1. Try Cache First!
	if cachedData, found := u.cache.Get(cacheKey); found {
		// Type assertion to convert generic interface{} back to *domain.Question
		if q, ok := cachedData.(*domain.Question); ok {
			return q, nil
		}
	}

	// 2. Cache Miss: Hit the Database
	q, err := u.questionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 3. Save to Cache for the next user!
	u.cache.Set(cacheKey, q)

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

// Upvote and DownVote for Quesions

func (u *questionUsecase) Upvote(userID string, questionID string) error {
	return u.handleVote(userID, questionID, 1) // 1 = Upvote
}

func (u *questionUsecase) Downvote(userID string, questionID string) error {
	return u.handleVote(userID, questionID, -1) // -1 = Downvote
}

// handleVote a private helper function to manage the complex math
func (u *questionUsecase) handleVote(userID string, questionID string, requestedValue int) error {
	// 1. Check if the user has already voted on this question
	existingVote, err := u.voteRepo.GetVote(userID, questionID)

	if err != nil { // NO VOTE EXISTS (New Vote)
		newVote := &domain.Vote{
			UserID:     userID,
			QuestionID: questionID,
			Value:      requestedValue,
		}

		// Add to Vote table
		if err := u.voteRepo.AddVote(newVote); err != nil {
			return err
		}

		// Update the cached total on the Question
		if requestedValue == 1 {
			return u.questionRepo.UpdateVoteCount(questionID, 1, 0)
		} else {
			return u.questionRepo.UpdateVoteCount(questionID, 0, 1)
		}
	}

	// 2. A VOTE ALREADY EXISTS

	// Scenario A: User clicked the exact same button (e.g. Upvoted an already upvoted post)
	// Action: "Toggle" or Cancel the vote.
	if existingVote.Value == requestedValue {
		if err := u.voteRepo.DeleteVote(userID, questionID); err != nil {
			return err
		}
		// Remove the vote from the totals
		if requestedValue == 1 {
			return u.questionRepo.UpdateVoteCount(questionID, -1, 0)
		} else {
			return u.questionRepo.UpdateVoteCount(questionID, 0, -1)
		}
	}

	// Scenario B: User is switching their vote (e.g. from Downvote to Upvote)
	// Action: Update the vote record, and adjust BOTH counters.
	existingVote.Value = requestedValue
	if err := u.voteRepo.UpdateVote(existingVote); err != nil {
		return err
	}

	if requestedValue == 1 {
		// Switching from Down to Up: Add 1 upvote, Subtract 1 downvote
		return u.questionRepo.UpdateVoteCount(questionID, 1, -1)
	} else {
		// Switching from Up to Down: Subtract 1 upvote, Add 1 downvote
		return u.questionRepo.UpdateVoteCount(questionID, -1, 1)
	}
}
