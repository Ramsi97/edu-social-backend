package middleware

import (
	"fmt"
	"net/http"

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

		userID, err := auth.ValidateToken(authHeader)

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
