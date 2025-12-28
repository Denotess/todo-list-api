package repository

import (
	"database/sql"

	"main.go/internal/models"
)

type TodoRepo struct {
	DB *sql.DB
}

func NewTodoRepo(db *sql.DB) *TodoRepo {
	return &TodoRepo{DB: db}
}

func (r *TodoRepo) GetTodos(userID int64, limit, offset int) ([]models.Todo, error) {
	rows, err := r.DB.Query(
		"SELECT id, user_id, title, content, is_done FROM todos WHERE user_id = ? ORDER BY id DESC LIMIT ? OFFSET ?;",
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]models.Todo, 0)
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.Id, &t.UserId, &t.Title, &t.Content, &t.IsDone); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *TodoRepo) InsertTodo(userId int64, todo models.CreateTodo) (int64, error) {
	res, err := r.DB.Exec(
		"INSERT INTO todos (user_id, title, content, is_done) VALUES (?, ?, ?, ?)",
		userId, todo.Title, todo.Content, todo.IsDone,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *TodoRepo) DeleteTodo(userId, todoId int64) error {
	_, err := r.DB.Exec("DELETE FROM todos WHERE id = ? AND user_id = ?;", todoId, userId)
	return err
}

func (r *TodoRepo) UpdateTodo(userId, todoId int64, title, content, isDone any) error {
	_, err := r.DB.Exec(
		"UPDATE todos SET title = COALESCE(?, title), content = COALESCE(?, content), is_done = COALESCE(?, is_done) WHERE user_id = ? AND id = ?;",
		title, content, isDone, userId, todoId,
	)
	return err
}
