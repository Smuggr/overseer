package api

import (
	"net/http"
	"os"
	"strings"

	"overseer/data/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)


func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrUnauthorized.Key})
			c.Abort()
			return
		}

		// "Bearer <token>"
		tokenString := strings.Split(header, " ")[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_TOKEN")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrInvalidToken.Key})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		login, ok := claims["login"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrInvalidTokenFormat.Key})
			c.Abort()
			return
		}

		c.Set("login", login)
		c.Next()
	}
}
