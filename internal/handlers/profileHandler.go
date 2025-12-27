package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"main.go/internal/db"
	"main.go/internal/models"
)

func GetTodos(ctx *gin.Context) {
	userIdStr, ok := ctx.Get("userId")
	if !ok {
		log.Println("no id in jwt")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}
	userId, err := strconv.ParseInt(userIdStr.(string), 10, 64)
	if err != nil {
		log.Println("error while parsing int")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	var query models.TodoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Println("error while binding query")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "wrong query parameters"})
		return
	}
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}
	if query.Offset < 0 {
		query.Offset = 0
	}

	rows, err := db.DB.Query("SELECT id, user_id, title, content, is_done FROM todos WHERE user_id = ? ORDER BY id DESC LIMIT ? OFFSET ?;", userId, query.Limit, query.Offset)
	if err != nil {
		log.Println("error while querying profile data from db")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer rows.Close()
	todos := make([]models.Todo, 0)
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.Id, &t.UserId, &t.Title, &t.Content, &t.IsDone); err != nil {
			fmt.Println("error while scanning rows in profileHandler")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
		log.Println("error while scanning the rows")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"todos": todos, "page": models.Page(&query), "limit": query.Limit, "offset": query.Offset})

}
func AddTodo(ctx *gin.Context) {
	userIdStr, ok := ctx.Get("userId")
	if !ok {
		log.Println("no id in jwt")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}
	userId, err := strconv.ParseInt(userIdStr.(string), 10, 64)
	if err != nil {
		log.Println("error while parsing int")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	var CreateTodo models.CreateTodo
	if err := ctx.ShouldBindJSON(&CreateTodo); err != nil {
		log.Println("json body error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect body parameters"})
		return
	}
	if strings.TrimSpace(CreateTodo.Content) == "" || strings.TrimSpace(CreateTodo.Title) == "" {
		log.Println("Content or title cannot be empty")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "content or title cannot be empty"})
		return
	}

	_, err = db.DB.Exec("INSERT INTO todos (user_id, title, content, is_done) VALUES (?, ?, ?, ?)", userId, CreateTodo.Title, CreateTodo.Content, CreateTodo.IsDone)
	if err != nil {
		log.Println("error while inserting todos")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"title": CreateTodo.Title, "content": CreateTodo.Content})
}

func DeleteTodo(ctx *gin.Context) {
	userIdStr, ok := ctx.Get("userId")
	if !ok {
		log.Println("no id in jwt")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}
	userId, err := strconv.ParseInt(userIdStr.(string), 10, 64)
	if err != nil {
		log.Println("error while parsing int")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	todoIdStr := ctx.Param("todoId") // JWT later
	todoId, err := strconv.ParseInt(todoIdStr, 10, 64)
	if err != nil {
		log.Println("error while parsing int")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	_, err = db.DB.Exec("DELETE FROM todos WHERE id = ? AND user_id = ?;", todoId, userId)
	if err != nil {
		log.Println("error while deleting todo")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func UpdateTodo(ctx *gin.Context) {
	// could improove with canonicalization
	userIdStr, ok := ctx.Get("userId")
	if !ok {
		log.Println("no id in jwt")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}
	userId, err := strconv.ParseInt(userIdStr.(string), 10, 64)
	if err != nil {
		log.Println("error while parsing int")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	todoIdStr := ctx.Param("todoId") // JWT later
	todoId, err := strconv.ParseInt(todoIdStr, 10, 64)
	if err != nil {
		log.Println("error while parsing int")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	var query models.UpdateTodo
	if err := ctx.ShouldBindJSON(&query); err != nil {
		fmt.Println("error while binding json")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}
	if query.Title == nil && query.Content == nil && query.IsDone == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}
	var title, content, isDone any

	if query.Title != nil {
		if strings.TrimSpace(*query.Title) == "" {
		} else {
			title = *query.Title
		}
	}

	if query.Content != nil {
		if strings.TrimSpace(*query.Content) == "" {
		} else {
			content = *query.Content
		}
	}

	if query.IsDone != nil {
		isDone = *query.IsDone
	}

	_, err = db.DB.Exec(
		"UPDATE todos SET title = COALESCE(?, title), content = COALESCE(?, content), is_done = COALESCE(?, is_done) WHERE user_id = ? AND id = ?;",
		title, content, isDone, userId, todoId,
	)
	if err != nil {
		log.Println("error while updating todo")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"title": query.Title, "content": query.Content})

}
