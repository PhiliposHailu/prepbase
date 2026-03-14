package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philipos/prepbase/delivery"
	"github.com/philipos/prepbase/domain"
	"github.com/philipos/prepbase/infrastructure"
)

func SetupRouter(userController *delivery.UserController, questionCtrl *delivery.QuestionController, cCtrl *delivery.CommentController, jwtSvc domain.JWTService) *gin.Engine {
	r := gin.Default()

	// Public Routes
	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)

	// Anyone can read questions!
	r.GET("/questions", questionCtrl.FetchAll)
	r.GET("/questions/:id", questionCtrl.GetByID)
	r.GET("/questions/:id/comments", cCtrl.GetByQuestionID)

	// Protected Routes (Require valid JWT)
	protected := r.Group("/users")
	protected.Use(infrastructure.AuthMiddleware(jwtSvc))
	{
		// Users
		protected.GET("/profile", userController.GetProfile)
		protected.PUT("/profile", userController.UpdateProfile)
		protected.DELETE("/:id", userController.DeleteUser)
		protected.PUT("/:id/promote", infrastructure.RoleMiddleware("admin"), userController.PromoteUser)

		// Questions
		protected.POST("/questions", questionCtrl.Create)
		protected.PUT("/questions/:id", questionCtrl.Update)
		protected.DELETE("/questions/:id", questionCtrl.Delete)
		protected.POST("/questions/:id/hint", questionCtrl.GenerateAIHint)

		// Upvote and Downvote Questions
		protected.POST("/questions/:id/upvote", questionCtrl.Upvote)
		protected.POST("/questions/:id/downvote", questionCtrl.Downvote)

		// Comments
		protected.POST("/questions/:id/comments", cCtrl.Create)
		protected.PUT("/comments/:comment_id", cCtrl.Update)
		protected.DELETE("/comments/:comment_id", cCtrl.Delete)

	}

	return r
}
