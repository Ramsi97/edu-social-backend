package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Ramsi97/edu-social-backend/pkg/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization format"})
			return
		}

		token := parts[1]

		userID, err := auth.ValidateToken(token)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		fmt.Println("UserID: " + userID)

		ctx.Set("user_id", userID)
		ctx.Next()
	}
}
