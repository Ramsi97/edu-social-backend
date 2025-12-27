package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	// Chat
	chatPostgres "github.com/Ramsi97/edu-social-backend/internal/chat/repository/postgres"
	chatUseCase "github.com/Ramsi97/edu-social-backend/internal/chat/use_case"
	chatHttp "github.com/Ramsi97/edu-social-backend/internal/chat/delivery/http"
	"github.com/Ramsi97/edu-social-backend/internal/chat/socket"

	// Auth feature
	authHttp "github.com/Ramsi97/edu-social-backend/internal/auth/delivery/http"
	authPostgres "github.com/Ramsi97/edu-social-backend/internal/auth/repository/postgres"
	authUseCase "github.com/Ramsi97/edu-social-backend/internal/auth/use_case"

	// Post feature
	postHttp "github.com/Ramsi97/edu-social-backend/internal/post/delivery/http"
	postPostgres "github.com/Ramsi97/edu-social-backend/internal/post/repository/postgres"
	postUseCase "github.com/Ramsi97/edu-social-backend/internal/post/use_case"

	// Like feature
	likeHttp "github.com/Ramsi97/edu-social-backend/internal/like/delivery/http"
	likePostgres "github.com/Ramsi97/edu-social-backend/internal/like/repository/postgres"
	likeUseCase "github.com/Ramsi97/edu-social-backend/internal/like/use_case"

	// Shared
	"github.com/Ramsi97/edu-social-backend/internal/middleware"
	cloud "github.com/Ramsi97/edu-social-backend/internal/shared/infrastructure"
	"github.com/Ramsi97/edu-social-backend/pkg/auth"
)

func main() {
	// -------------------
	// Load configuration
	// -------------------
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not found in env")
	}
	auth.SetJWTSecret(jwtSecret)

	dbHost := "localhost"
	dbPort := 5432
	dbUser := "ramsi"
	dbPassword := "RAMSIDB"
	dbName := "edu_social_db"

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	// -------------------
	// Initialize Cloudinary
	// -------------------
	cldInstance, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Fatal("Failed to init Cloudinary")
	}
	mediaUploader := cloud.NewCloudinaryUploader(cldInstance)

	// -------------------
	// Initialize Repositories
	// -------------------
	userRepo := authPostgres.NewUserRepository(db)
	postRepo := postPostgres.NewPostRepository(db)
	likeRepo := likePostgres.NewLikeRepository(db)
	chatRepo := chatPostgres.NewChatRepository(db)

	// -------------------
	// Initialize Use Cases
	// -------------------
	authUC := authUseCase.NewAuthUseCase(userRepo, mediaUploader)
	postUC := postUseCase.NewPostUseCase(postRepo)
	likeUC := likeUseCase.NewLikeUseCase(likeRepo)
	chatUC := chatUseCase.NewChatUseCase(chatRepo)

	// -------------------
	// Initialize Router
	// -------------------
	router := gin.Default()

	// Health check
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// -------------------
	// Start Socket.IO
	// -------------------
	socketServer, err := socket.StartSocketServer()
	if err != nil {
		log.Fatal("Socket.IO server failed:", err)
	}

	// Attach Socket.IO to Gin
	router.GET("/socket.io/*any", gin.WrapH(socketServer))
	router.POST("/socket.io/*any", gin.WrapH(socketServer))

	// -------------------
	// API Groups
	// -------------------
	api := router.Group("/api/v1")
	authGroup := api.Group("/auth")
	postGroup := api.Group("/post")
	postGroup.Use(middleware.AuthMiddleWare())
	likeGroup := api.Group("/like")
	likeGroup.Use(middleware.AuthMiddleWare())
	chatGroup := api.Group("/chat")
	chatGroup.Use(middleware.AuthMiddleWare())

	// -------------------
	// Attach Handlers
	// -------------------
	authHttp.NewAuthHandler(authGroup, authUC)
	postHttp.NewPostHandler(postGroup, postUC, mediaUploader)
	likeHttp.NewLikeHandler(likeGroup, likeUC)
	chatHttp.NewChatHandler(chatGroup, chatUC)

	// -------------------
	// Run server
	// -------------------
	router.Run(":8080")
}
