package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"

	_ "github.com/lib/pq"
	"github.com/zishang520/socket.io/v2/socket"
	_ "github.com/jackc/pgx/v5/stdlib"

	// Chat
	chatHttp "github.com/Ramsi97/edu-social-backend/internal/chat/delivery/http"
	chatPostgres "github.com/Ramsi97/edu-social-backend/internal/chat/repository/postgres"
	chatSocket "github.com/Ramsi97/edu-social-backend/internal/chat/socket"
	chatUseCase "github.com/Ramsi97/edu-social-backend/internal/chat/use_case"

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

	// Comment Feature
	commentHttp "github.com/Ramsi97/edu-social-backend/internal/comment/delivery/http"
	commentPostgres "github.com/Ramsi97/edu-social-backend/internal/comment/repository/postgres"
	commentUseCase "github.com/Ramsi97/edu-social-backend/internal/comment/use_case"

	// Group Chat Feature
	groupHttp "github.com/Ramsi97/edu-social-backend/internal/group/delivery/http"
	groupSocket "github.com/Ramsi97/edu-social-backend/internal/group/delivery/socket"
	groupPostgres "github.com/Ramsi97/edu-social-backend/internal/group/repository/postgres"
	groupUseCase "github.com/Ramsi97/edu-social-backend/internal/group/use_case"

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

	dbHost := os.Getenv("PGHOST")
	dbUser := os.Getenv("PGUSER")
	dbPassword := os.Getenv("PGPASSWORD")
	dbName := os.Getenv("PGDATABASE")
	dbSSLMode := os.Getenv("PGSSLMODE")
	port := 5432

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbHost, port, dbUser, dbPassword, dbName, dbSSLMode,
	)

	db, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Connected to Neon PostgreSQL")

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
	commentRepo := commentPostgres.NewCommentRepository(db)
	chatRepo := chatPostgres.NewChatRepository(db)
	groupchatRepo := groupPostgres.NewGroupChatRepo(db)

	// -------------------
	// Initialize Use Cases
	// -------------------
	authUC := authUseCase.NewAuthUseCase(userRepo, mediaUploader)
	postUC := postUseCase.NewPostUseCase(postRepo)
	likeUC := likeUseCase.NewLikeUseCase(likeRepo)
	commentUC := commentUseCase.NewCommentUseCase(commentRepo)
	chatUC := chatUseCase.NewChatUseCase(chatRepo)
	groupchatUC := groupUseCase.NewGroupChatUseCase(groupchatRepo)

	// -------------------
	// Initialize Router
	// -------------------
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	// Health check
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// ----------------------------------
	// initialize model Socket.IO Server
	// ----------------------------------
	io := socket.NewServer(nil, nil)

	// Register chat (1–1)
	chatSocketHandler := chatSocket.NewSocketHandler(io, chatUC)
	chatSocketHandler.RegisterMiddleWare()
	chatSocketHandler.RegisterEvents()

	// Register group chat
	groupChatSocketHandler := groupSocket.NewSocketHandler(io, groupchatUC)
	groupChatSocketHandler.RegisterMiddleWare()
	groupChatSocketHandler.RegisterEvents()

	router.GET("/socket.io/*any", gin.WrapH(io.ServeHandler(nil)))
	router.POST("/socket.io/*any", gin.WrapH(io.ServeHandler(nil)))

	// -------------------
	// API Groups
	// -------------------
	api := router.Group("/api/v1")
	authGroup := api.Group("/auth")
	postGroup := api.Group("/posts")
	postGroup.Use(middleware.AuthMiddleWare())
	likeGroup := api.Group("/like")
	likeGroup.Use(middleware.AuthMiddleWare())
	commentGroup := api.Group("/comment")
	commentGroup.Use(middleware.AuthMiddleWare())
	likeGroup.Use(middleware.AuthMiddleWare())
	chatGroup := api.Group("/chat")
	chatGroup.Use(middleware.AuthMiddleWare())
	groupApiGroup := api.Group("/group")
	groupApiGroup.Use(middleware.AuthMiddleWare())

	// -------------------
	// Attach Handlers
	// -------------------
	authHttp.NewAuthHandler(authGroup, authUC)
	postHttp.NewPostHandler(postGroup, postUC, mediaUploader)
	likeHttp.NewLikeHandler(likeGroup, likeUC)
	commentHttp.NewCommentHandler(commentGroup, commentUC)
	chatHttp.NewChatHandler(chatGroup, chatUC)
	groupHttp.NewGroupHandler(groupchatUC, groupApiGroup)

	// -------------------
	// Run server
	// -------------------
	router.Run(":8080")
}
