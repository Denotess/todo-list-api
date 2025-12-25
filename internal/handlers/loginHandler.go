package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/internal/db"
	"main.go/internal/helpers"
	"main.go/internal/models"
)

func Login(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	exists, err := helpers.CheckDuplicateUser(user.Name)
	if err != nil {
		log.Println("error while checking if user exists")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if !exists {
		log.Println("user does not yet exist")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user does not exist"})
		return
	}
	var hashedPass string
	if err := db.DB.QueryRow("SELECT password_hash FROM users WHERE name = ?", user.Name).Scan(&hashedPass); err != nil {
		log.Println("error while selectin password from db")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	passwordMatch, err := helpers.CheckPasswordHash(user.Password, hashedPass)
	if err != nil {
		log.Println("error while checking if passwords match")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if !passwordMatch {
		log.Println("Passwords dont match")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect password"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "succesfull login"}) // returns a token here
}
