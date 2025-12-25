package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/internal/db"
	"main.go/internal/helpers"
)

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func Register(ctx *gin.Context) {
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	if user.Name == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name and password required"})
		return
	}
	isDupe, err := helpers.CheckDuplicateUser(user.Name)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if isDupe {
		log.Println("user already exists")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username already taken"})
		return
	}
	passwordHash, err := helpers.HashPassword(user.Password)
	if err != nil {
		log.Println("error while hashing pass: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error while hashing pass"})
		return
	}
	_, err = db.DB.Exec("INSERT INTO users (name, password_hash) values (?, ?)", user.Name, passwordHash)
	if err != nil {
		log.Println("register insert error")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting user"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "ok"})

}
