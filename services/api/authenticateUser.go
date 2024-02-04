package api

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"overseer/data/errors"
	"overseer/data/models"
	"overseer/services/database"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)


func createToken(login string) (string, error) {
	jwt_token_lifespan, err := strconv.Atoi(os.Getenv("API_JWT_TOKEN_LIFESPAN_MINUTES"))
	if err != nil {
		return "", errors.ErrCreatingToken
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(time.Duration(jwt_token_lifespan) * time.Minute).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("SECRET_TOKEN")))
}

func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrUnauthorized.Error()})
			c.Abort()
			return
		}

		// "Bearer <token>"
		tokenString := strings.Split(header, " ")[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_TOKEN")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrInvalidToken.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthenticateUser(c *gin.Context) {
    var user models.User
    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrInvalidRequestPayload.Error()})
        return
    }

	// User not found? Or just invalid credentials?
    var existingUser models.User
    if database.DB.Where("login = ?", user.Login).First(&existingUser).Error != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrInvalidCredentials.Error()})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrInvalidCredentials.Error()})
        return
    }

    tokenString, err := createToken(user.Login)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}