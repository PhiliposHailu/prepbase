package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/philipos/prepbase/delivery"
	"github.com/philipos/prepbase/infrastructure"
	"github.com/philipos/prepbase/repository"
	"github.com/philipos/prepbase/router"
	"github.com/philipos/prepbase/usecase"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Proceeding with system environment variables.")
	}

	port := os.Getenv("PORT")
	dbURI := os.Getenv("DB_URI")
	jwtSecret := os.Getenv("JWT_ACCESS_SECRET")

	if port == "" {
		log.Fatal("ERROR: PORT environment variable is required")
	}
	if dbURI == "" {
		log.Fatal("ERROR: DB_URI environment variable is required")
	}
	if jwtSecret == "" {
		log.Fatal("ERROR: JWT_ACCESS_SECRET environment variable is required")
	}

	db := infrastructure.ConnectDB()

	// Services
	pwdSvc := infrastructure.NewPasswordService()
	jwtSvc := infrastructure.NewJWTService()

	// Repositories
	userRepo := repository.NewUserRepository(db, "users")
	voteRepo := repository.NewVoteRepository(db, "votes")
	questionRepo := repository.NewQuestionRepository(db, "questions")
	commentRepo := repository.NewCommentRepository(db, "comments") 
	cacheSvc := infrastructure.NewMemoryCache()
	aiSvc := infrastructure.NewAIService()


	// Usecases
	userUsecase := usecase.NewUserUsecase(userRepo, pwdSvc, jwtSvc)
	questionUsecase := usecase.NewQuestionUsecase(questionRepo, voteRepo, cacheSvc, aiSvc) 
	commentUsecase := usecase.NewCommentUsecase(commentRepo)

	// Controllers
	userController := delivery.NewUserController(userUsecase)
	questionController := delivery.NewQuestionController(questionUsecase) 
	commentController := delivery.NewCommentController(commentUsecase)

	// Router
	r := router.SetupRouter(userController, questionController, commentController, jwtSvc)

	log.Println("🚀 Server running on port 8080 ??? ;)")
	r.Run(":8000")

}
