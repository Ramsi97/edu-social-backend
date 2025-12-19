package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Ramsi97/edu-social-backend/internal/auth/delivery/https"
	"github.com/Ramsi97/edu-social-backend/internal/auth/repository/postgres"
	usecase "github.com/Ramsi97/edu-social-backend/internal/auth/use_case"
	"github.com/gin-gonic/gin"
)

func main(){
	
	host := "localhost"
	port := 5432
	user := "ramsi"
	password := "RAMSIDB"
	dbname := "edu_socail_db"

	psqlinfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", psqlinfo)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	userRepo := postgres.NewUserRepository(db)
	authUseCase := usecase.NewAuthUseCase(userRepo)
	router := gin.Default()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	api := router.Group("/api/v1")

	authGroup := api.Group("/auth")
	
	https.NewAuthHandler(authGroup, authUseCase)

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	router.Run(":", serverPort)

}