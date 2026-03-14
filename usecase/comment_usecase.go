package usecase

import (
	"errors"
	"strings"

	"github.com/philipos/prepbase/domain"
)

type commentUsecase struct {
	commentRepo domain.CommentRepository
}

func NewCommentUsecase(cRepo domain.CommentRepository) domain.CommentUsecase {
	return &commentUsecase{commentRepo: cRepo}
}

func (u *commentUsecase) Create(c *domain.Comment) error {
	if strings.TrimSpace(c.Content) == "" {
		return errors.New("comment content cannot be empty")
	}
	if c.QuestionID == "" || c.AuthorID == "" {
		return errors.New("invalid comment associations")
	}
	return u.commentRepo.Create(c)
}

func (u *commentUsecase) GetByQuestionID(questionID string) ([]domain.Comment, error) {
	return u.commentRepo.GetByQuestionID(questionID)
}

func (u *commentUsecase) Update(id string, authorID string, content string) error {
	if strings.TrimSpace(content) == "" {
		return errors.New("content cannot be empty")
	}

	existingComment, err := u.commentRepo.GetByID(id)
	if err != nil {
		return errors.New("comment not found")
	}

	// AUTH: Only the author can update
	if existingComment.AuthorID != authorID {
		return errors.New("unauthorized: you can only update your own comments")
	}

	existingComment.Content = content
	return u.commentRepo.Update(existingComment)
}

func (u *commentUsecase) Delete(id string, authorID string, userRole string) error {
	existingComment, err := u.commentRepo.GetByID(id)
	if err != nil {
		return errors.New("comment not found")
	}

	// AUTH: Author OR Admin can delete
	if existingComment.AuthorID != authorID && userRole != "admin" {
		return errors.New("unauthorized: you do not have permission to delete this comment")
	}

	return u.commentRepo.Delete(id)
}