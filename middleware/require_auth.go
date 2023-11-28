package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sidx04/go-jwt/initialisers"
	"github.com/sidx04/go-jwt/models"
)

func RequireAuth(c *gin.Context) {
	/* get cookie off request */
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	/* decode and validate */
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		return jwtSecret, nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		/* check expiry */
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		/* find user with token `sub` */
		var user *models.User
		initialisers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		/* attach to request */
		c.Set("user", user)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	/* continue	*/
	c.Next()
}
