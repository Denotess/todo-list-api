package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"main.go/internal/db"
	"main.go/internal/handlers"
	"main.go/internal/middleware"
)

func main() {
	godotenv.Load()
	db.InitDB()
	defer db.DB.Close()
	router := gin.Default()
	authorized := router.Group("/users")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("todos", handlers.GetTodos)
		authorized.POST("todos", handlers.AddTodo)
		authorized.DELETE("todos/:todoId", handlers.DeleteTodo)
		authorized.PUT("todos/:todoId", handlers.UpdateTodo)
	}
	router.GET("/ping", handlers.Ping)
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)
	router.Run()
}
