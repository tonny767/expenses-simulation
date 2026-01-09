package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"backend/db"
	"backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.First(&user, "email = ?", input.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password doesn't match"})
		return
	}

	// Generate JWT token
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 1).Unix(), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Get environment
	maxAge := 60 * 60 * 24 * 1

	c.SetCookie(
		"token",
		tokenString,
		maxAge, // 24 hours
		"/",
		"",
		false,
		true, // httpOnly
	)

	// Set user info cookies (accessible by frontend)
	c.Writer.Header().Add("Set-Cookie", fmt.Sprintf("user_id=%d; Path=/; Max-Age=%d; SameSite=Lax", user.ID, maxAge))
	c.Writer.Header().Add("Set-Cookie", fmt.Sprintf("role=%s; Path=/; Max-Age=%d; SameSite=Lax", user.Role, maxAge))
	c.Writer.Header().Add("Set-Cookie", fmt.Sprintf("user_name=%s; Path=/; Max-Age=%d; SameSite=Lax", user.Name, maxAge))

	// Also return in response body for immediate use
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"id":    user.ID,
		"role":  user.Role,
		"name":  user.Name,
		"email": user.Email,
	})
}
