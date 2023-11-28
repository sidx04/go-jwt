package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sidx04/go-jwt/initialisers"
	"github.com/sidx04/go-jwt/models"
	"golang.org/x/crypto/bcrypt"
)

func UserSignUp(c *gin.Context) {
	/* get the email and password off req body */
	var requestBody struct {
		Email    string
		Password string
	}

	if c.Bind(&requestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body...",
		})
		return
	}

	/* hash pass */
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not hash password...",
		})
		return
	}

	/* create the user */
	user := models.User{
		Email:    requestBody.Email,
		Password: string(hashedPassword),
	}
	res := initialisers.DB.Create(&user)

	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user...",
		})
		return
	}

	/* respond */
	c.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
}

func LoginUser(c *gin.Context) {
	/* get email and password off req body */
	var requestBody struct {
		Email    string
		Password string
	}

	if c.Bind(&requestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body...",
		})
		return
	}

	/* query requested user */
	var user models.User
	initialisers.DB.First(&user, "email = ?", requestBody.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email...",
		})
		return
	}

	/* compare password with saved user hashed password */
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password...",
		})
		return
	}

	/* generate JSON web token */
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 10).Unix(), // expires in 10 days
	})

	// sign and get the complete encoded token as a string using secret key
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token...",
		})
		return
	}

	/* respond with a cookie */
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*10, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

func ValidateToken(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"user":   user,
		"status": "logged in",
	})
}
