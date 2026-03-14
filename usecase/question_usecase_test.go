package usecase_test

import (
	"errors"
	"testing"

	"github.com/philipos/prepbase/domain"
	"github.com/philipos/prepbase/mocks"
	"github.com/philipos/prepbase/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQuestionUsecase_Upvote(t *testing.T) {
	// Table-Driven Test Setup
	tests := []struct {
		name          string
		userID        string
		questionID    string
		mockBehavior  func(qRepo *mocks.QuestionRepository, vRepo *mocks.VoteRepository)
		expectedError string
	}{
		{
			name:       "Success - First Time Upvote",
			userID:     "user1",
			questionID: "q1",
			mockBehavior: func(qRepo *mocks.QuestionRepository, vRepo *mocks.VoteRepository) {
				// 1. Usecase checks if vote exists -> It doesn't
				vRepo.On("GetVote", "user1", "q1").Return(nil, errors.New("not found"))
				// 2. Usecase adds vote
				vRepo.On("AddVote", mock.AnythingOfType("*domain.Vote")).Return(nil)
				// 3. Usecase increments upvote count by 1, downvote by 0
				qRepo.On("UpdateVoteCount", "q1", 1, 0).Return(nil)
			},
			expectedError: "",
		},
		{
			name:       "Success - Toggle Upvote Off",
			userID:     "user1",
			questionID: "q1",
			mockBehavior: func(qRepo *mocks.QuestionRepository, vRepo *mocks.VoteRepository) {
				// 1. Usecase checks if vote exists -> It does, and it's already an Upvote (1)
				existingVote := &domain.Vote{UserID: "user1", QuestionID: "q1", Value: 1}
				vRepo.On("GetVote", "user1", "q1").Return(existingVote, nil)
				
				// 2. Usecase deletes the vote
				vRepo.On("DeleteVote", "user1", "q1").Return(nil)
				
				// 3. Usecase decrements the upvote count
				qRepo.On("UpdateVoteCount", "q1", -1, 0).Return(nil)
			},
			expectedError: "",
		},
		{
			name:       "Success - Switch Downvote to Upvote",
			userID:     "user1",
			questionID: "q1",
			mockBehavior: func(qRepo *mocks.QuestionRepository, vRepo *mocks.VoteRepository) {
				// 1. Existing vote is a Downvote (-1)
				existingVote := &domain.Vote{UserID: "user1", QuestionID: "q1", Value: -1}
				vRepo.On("GetVote", "user1", "q1").Return(existingVote, nil)
				
				// 2. Usecase updates the vote
				vRepo.On("UpdateVote", mock.AnythingOfType("*domain.Vote")).Return(nil)
				
				// 3. Usecase increments Upvote (+1) AND decrements Downvote (-1)
				qRepo.On("UpdateVoteCount", "q1", 1, -1).Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup Mocks
			mockQRepo := new(mocks.QuestionRepository)
			mockVRepo := new(mocks.VoteRepository)
			mockCache := new(mocks.CacheService) // We don't use cache in upvote, but constructor needs it
			mockAI := new(mocks.AIService)       // Same here

			// Apply the scripted behavior
			tc.mockBehavior(mockQRepo, mockVRepo)

			// Initialize the Usecase
			u := usecase.NewQuestionUsecase(mockQRepo, mockVRepo, mockCache, mockAI)

			// Execute
			err := u.Upvote(tc.userID, tc.questionID)

			// Assert Error
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			}

			// Verify Mocks
			mockQRepo.AssertExpectations(t)
			mockVRepo.AssertExpectations(t)
		})
	}
}