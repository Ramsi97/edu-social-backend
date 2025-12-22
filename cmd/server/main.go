package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Ramsi97/edu-social-backend/internal/auth/delivery/https"
	authPostgres "github.com/Ramsi97/edu-social-backend/internal/auth/repository/postgres"
	authusecase "github.com/Ramsi97/edu-social-backend/internal/auth/use_case"
	"github.com/Ramsi97/edu-social-backend/internal/middleware"
	posthandler "github.com/Ramsi97/edu-social-backend/internal/post/delivery/http"
	postPostgres "github.com/Ramsi97/edu-social-backend/internal/post/repository/postgres"
	postUseCase "github.com/Ramsi97/edu-social-backend/internal/post/use_case"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	host := "localhost"
	port := 5432
	user := "ramsi"
	password := "RAMSIDB"
	dbname := "edu_social_db"

	psqlinfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", psqlinfo)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	userRepo := authPostgres.NewUserRepository(db)
	authUseCase := authusecase.NewAuthUseCase(userRepo)
	router := gin.Default()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	postRepo := postPostgres.NewPostRepository(db)
	postUseCase := postUseCase.NewPostUseCase(postRepo)

	api := router.Group("/api/v1")

	authGroup := api.Group("/auth")
	postGroup := api.Group("/post")
	postGroup.Use(middleware.AuthMiddleWare())
	https.NewAuthHandler(authGroup, authUseCase)
	posthandler.NewPostHandler(postGroup, postUseCase)

	router.Run()

}
