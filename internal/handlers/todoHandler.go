package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/internal/models"
	"main.go/internal/service"
)

type TodoHandler struct {
	Service *service.TodoService
}

func NewTodoHandler(service *service.TodoService) *TodoHandler {
	return &TodoHandler{Service: service}
}

func (h *TodoHandler) GetTodos(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}

	var query models.TodoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "wrong query parameters"})
		return
	}

	todos, query, err := h.Service.GetTodos(userID, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"todos":  todos,
		"page":   models.Page(&query),
		"limit":  query.Limit,
		"offset": query.Offset,
	})
}

func (h *TodoHandler) AddTodo(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}

	var payload models.CreateTodo
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect body parameters"})
		return
	}

	todoID, err := h.Service.AddTodo(userID, payload)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTodo) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "content or title cannot be empty"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":      todoID,
		"title":   payload.Title,
		"content": payload.Content,
	})
}

func (h *TodoHandler) DeleteTodo(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}
	todoID, ok := getTodoID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}

	if err := h.Service.DeleteTodo(userID, todoID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *TodoHandler) UpdateTodo(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}
	todoID, ok := getTodoID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}

	var payload models.UpdateTodo
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}

	if err := h.Service.UpdateTodo(userID, todoID, payload); err != nil {
		if errors.Is(err, service.ErrNoUpdateFields) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func getUserID(ctx *gin.Context) (int64, bool) {
	raw, ok := ctx.Get("userId")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return 0, false
	}
	userIDStr, ok := raw.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return 0, false
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return 0, false
	}
	return userID, true
}

func getTodoID(ctx *gin.Context) (int64, bool) {
	todoIDStr := ctx.Param("todoId")
	todoID, err := strconv.ParseInt(todoIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return 0, false
	}
	return todoID, true
}
