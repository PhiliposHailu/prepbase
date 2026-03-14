package delivery

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/philipos/prepbase/domain"
)

type QuestionController struct {
	qUsecase domain.QuestionUsecase
}

func NewQuestionController(qu domain.QuestionUsecase) *QuestionController {
	return &QuestionController{
		qUsecase: qu,
	}
}

func (qc *QuestionController) GenerateAIHint(c *gin.Context) {
	id := c.Param("id")

	hint, err := qc.qUsecase.GenerateAIHint(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hint": hint})
}

func (qc *QuestionController) Create(c *gin.Context) {
	var q domain.Question
	if err := c.BindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Security: Force the AuthorID to be the currently logged-in user!
	// (Never trust the client to send their own ID in the JSON body)
	q.AuthorID = c.GetString("user_id")

	if err := qc.qUsecase.Create(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Question posted successfully", "question": q})
}

func (qc *QuestionController) GetByID(c *gin.Context) {
	id := c.Param("id")

	q, err := qc.qUsecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, q)
}

func (qc *QuestionController) FetchAll(c *gin.Context) {
	// Read pagination parameters from the URL query string (e.g., ?limit=10&cursor=timestamp)
	cursor := c.Query("cursor")
	limitStr := c.Query("limit")

	limit := 10 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	questions, err := qc.qUsecase.FetchAll(cursor, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func (qc *QuestionController) Update(c *gin.Context) {
	id := c.Param("id")
	authorID := c.GetString("user_id")

	var q domain.Question
	if err := c.BindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := qc.qUsecase.Update(id, authorID, &q); err != nil {
		// If error contains "unauthorized", send 403. Otherwise 400.
		status := http.StatusBadRequest
		if err.Error() == "unauthorized: you can only update your own questions" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question updated successfully"})
}

func (qc *QuestionController) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	userRole := c.GetString("role")

	if err := qc.qUsecase.Delete(id, userID, userRole); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "unauthorized: you do not have permission to delete this question" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}

func (qc *QuestionController) Upvote(c *gin.Context) {
	questionID := c.Param("id")
	userID := c.GetString("user_id")

	if err := qc.qUsecase.Upvote(userID, questionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upvote"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vote registered"})
}

func (qc *QuestionController) Downvote(c *gin.Context) {
	questionID := c.Param("id")
	userID := c.GetString("user_id")

	if err := qc.qUsecase.Downvote(userID, questionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to downvote"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vote registered"})
}
