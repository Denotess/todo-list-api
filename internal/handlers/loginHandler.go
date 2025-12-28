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
	if user.Name == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name and password required"})
		return
	}
	exists, err := helpers.CheckDuplicateUser(user.Name)
	if err != nil {
		log.Println("error while checking if user exists", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if !exists {
		log.Println("user does not yet exist")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user does not exist"})
		return
	}
	var hashedPass string
	if err := db.DB.QueryRow("SELECT id, password_hash FROM users WHERE name = ?", user.Name).Scan(&user.Id, &hashedPass); err != nil {
		log.Println("error while selectin password from db", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	passwordMatch, err := helpers.CheckPasswordHash(user.Password, hashedPass)
	if err != nil {
		log.Println("error while checking if passwords match", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if !passwordMatch {
		log.Println("Passwords dont match", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect password"})
		return
	}
	token, err := helpers.CreateToken(&user)
	if err != nil {
		log.Println("error while creating token", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "succesfull login", "token": token}) // returns a token here
}
