package service

import (
	"errors"
	"strings"

	"main.go/internal/models"
	"main.go/internal/repository"
)

var ErrNoUpdateFields = errors.New("no fields to update")
var ErrInvalidTodo = errors.New("tile and content required")

type TodoService struct {
	Repo *repository.TodoRepo
}

func NewTodoService(repo *repository.TodoRepo) *TodoService {
	return &TodoService{Repo: repo}
}
func (s *TodoService) GetTodos(userId int64, query models.TodoQuery) ([]models.Todo, models.TodoQuery, error) {
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}
	if query.Offset < 0 {
		query.Offset = 0
	}
	todos, err := s.Repo.GetTodos(userId, query.Limit, query.Offset)
	return todos, query, err
}

func (s *TodoService) AddTodo(userID int64, todo models.CreateTodo) (int64, error) {
	if strings.TrimSpace(todo.Title) == "" || strings.TrimSpace(todo.Content) == "" {
		return 0, ErrInvalidTodo
	}
	return s.Repo.InsertTodo(userID, todo)
}

func (s *TodoService) DeleteTodo(userID, todoID int64) error {
	return s.Repo.DeleteTodo(userID, todoID)
}

func (s *TodoService) UpdateTodo(userID, todoID int64, update models.UpdateTodo) error {
	if update.Title == nil && update.Content == nil && update.IsDone == nil {
		return ErrNoUpdateFields
	}

	var title, content, isDone any

	if update.Title != nil && strings.TrimSpace(*update.Title) != "" {
		title = *update.Title
	}
	if update.Content != nil && strings.TrimSpace(*update.Content) != "" {
		content = *update.Content
	}
	if update.IsDone != nil {
		isDone = *update.IsDone
	}

	return s.Repo.UpdateTodo(userID, todoID, title, content, isDone)
}
