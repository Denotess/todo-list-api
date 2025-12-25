package main

import (
	"github.com/gin-gonic/gin"
	"main.go/internal/db"
	"main.go/internal/handlers"
)

func main() {
	db.InitDB()
	defer db.DB.Close()
	router := gin.Default()
	router.GET("/ping", handlers.Ping)
	router.POST("/register", handlers.Register)
	router.Run()
}
