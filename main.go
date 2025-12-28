package main

import (
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"main.go/db"
	"main.go/handlers"
	"main.go/middleware"
)

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

func main() {
	godotenv.Load()
	db.InitDB()
	defer db.DB.Close()
	router := gin.Default()
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Second,
		Limit: 5,
	})
	rateLimiter := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})
	api := router.Group("/api")
	api.Use(rateLimiter)
	{

		authorized := api.Group("/users")
		authorized.Use(middleware.AuthMiddleware())
		{
			authorized.GET("/todos", handlers.GetTodos)
			authorized.POST("/todos", handlers.AddTodo)
			authorized.DELETE("/todos/:todoId", handlers.DeleteTodo)
			authorized.PUT("/todos/:todoId", handlers.UpdateTodo)
		}
		api.GET("/ping", handlers.Ping)
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
	}
	router.Run()
}
