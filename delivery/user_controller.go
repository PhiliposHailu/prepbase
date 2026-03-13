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

func (uc *UserController) Register(c *gin.Context) {
	var user domain.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	if err := uc.userUsecase.Register(&user); err != nil {
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