package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/internal/db"
	"main.go/internal/models"
)

func GetProfile(ctx *gin.Context) {
	var query models.TodoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Println("error while binding query")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "wrong query parameters"})
		return
	}
	rows, err := db.DB.Query("SELECT id, user_id, title, content, is_done FROM todos WHERE user_id = ? ORDER BY id DESC LIMIT ? OFFSET ?;", &todoData.UserId, &query.Limit, &query.Offset)
	if err != nil {
		log.Println("error while querying profile data from db")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer rows.Close()
	todos := make([]models.Todo, 0)
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.Id, &t.UserId, &t.Title, &t.Content, t.IsDone); err != nil {
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
