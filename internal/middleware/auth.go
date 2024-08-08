package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"ticket_goroutine/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]
		log.Info().Msg(tokenString)
		claims := &domain.Claims{}
		secret := []byte(os.Getenv("SECRET"))

		token, err := jwt.ParseWithClaims(tokenString, claims,
			func (token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("Error")
				}
				return secret, nil
			},
		)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, fmt.Sprintf("Welcome %s", claims.Name))
		c.Next()
	}
}