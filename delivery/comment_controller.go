package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/philipos/prepbase/domain"
)

type CommentController struct {
	cUsecase domain.CommentUsecase
}

func NewCommentController(cu domain.CommentUsecase) *CommentController {
	return &CommentController{cUsecase: cu}
}

func (cc *CommentController) Create(c *gin.Context) {
	var comment domain.Comment
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	comment.AuthorID = c.GetString("user_id")
	// The QuestionID will be passed in the URL (e.g. /questions/:id/comments)
	comment.QuestionID = c.Param("question_id") 

	if err := cc.cUsecase.Create(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Comment posted", "comment": comment})
}

func (cc *CommentController) GetByQuestionID(c *gin.Context) {
	qID := c.Param("question_id")
	comments, err := cc.cUsecase.GetByQuestionID(qID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (cc *CommentController) Update(c *gin.Context) {
	id := c.Param("comment_id")
	authorID := c.GetString("user_id")

	var input struct {
		Content string `json:"content"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := cc.cUsecase.Update(id, authorID, input.Content); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment updated"})
}

func (cc *CommentController) Delete(c *gin.Context) {
	id := c.Param("comment_id")
	userID := c.GetString("user_id")
	userRole := c.GetString("role")

	if err := cc.cUsecase.Delete(id, userID, userRole); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}