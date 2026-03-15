package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/philipos/prepbase/domain"
)

type UserController struct {
	userUsecase domain.UserUsecase
}

func NewUserController(u domain.UserUsecase) *UserController {
	return &UserController{
		userUsecase: u,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Bio      string `json:"bio"`
	Role     string `json:"role"`
}

func (uc *UserController) Register(c *gin.Context) {
	var user RegisterRequest
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	domainUser := &domain.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Bio:      user.Bio,
		Role:     user.Role,
	}
	if err := uc.userUsecase.Register(domainUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

func (uc *UserController) Login(c *gin.Context) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	accessToken, refreshToken, err := uc.userUsecase.Login(creds.Email, creds.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (uc *UserController) GetProfile(c *gin.Context) {
	// Grab the ID of the currently logged-in user from the context
	userID := c.GetString("user_id")

	user, err := uc.userUsecase.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uc *UserController) UpdateProfile(c *gin.Context) {
	var user domain.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user.ID = c.GetString("user_id")

	if err := uc.userUsecase.UpdateProfile(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func (uc *UserController) PromoteUser(c *gin.Context) {
	targetID := c.Param("id")
	adminID := c.GetString("user_id")

	if err := uc.userUsecase.PromoteUser(adminID, targetID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin successfully"})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	targetID := c.Param("id")
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")

	err := uc.userUsecase.DeleteUser(actorID, actorRole, targetID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deactivated successfully"})
}

func (uc *UserController) RefreshToken(c *gin.Context) {
	var body struct { RefreshToken string `json:"refresh_token"` }
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	newAccess, err := uc.userUsecase.RefreshToken(body.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"access_token": newAccess})
}

func (uc *UserController) ForgotPassword(c *gin.Context) {
	var body struct { Email string `json:"email"` }
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := uc.userUsecase.ForgotPassword(body.Email); err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "If the email exists, a reset link has been sent."})
}