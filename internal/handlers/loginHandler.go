package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/internal/models"
)

func Login(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
	}
}
