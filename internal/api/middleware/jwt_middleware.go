package middleware

import (
	"net/http"
	"strings"

	"github.com/Masha003/Golang-internship.git/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JwtAuth(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := extract(ctx, secret)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		id, err := auth.ExtractId(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Set("user_id", id)
		ctx.Next()
	}
}

func extract(ctx *gin.Context, secret string) (*jwt.Token, error) {
	tokenString := ctx.Query("token")
	if tokenString == "" {
		bearerToken := ctx.Request.Header.Get("Authorization")
		if len(strings.Split(bearerToken, " ")) == 2 {
			tokenString = strings.Split(bearerToken, " ")[1]
		}
	}

	return auth.Validate(tokenString, secret)
}
