package repository

import (
	"database/sql"

	"main.go/internal/models"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CheckDuplicateUser(name string) (bool, error) {
	var existing string
	err := r.DB.QueryRow("SELECT name FROM users WHERE name = ?", name).Scan(&existing)
	if err == sql.ErrNoRows {
		return false, err
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserRepo) GetUserByName(name string) (models.User, error) {
	var user models.User
	if err := r.DB.QueryRow("SELECT id, name, password_hash FROM users WHERE name = ?", name).Scan(&user.Id, &user.Name, &user.PasswordHash); err != nil {
		return models.User{}, err
	}
	return user, nil
}
func (r *UserRepo) InsertUser(name, passwordHash string) (int64, error) {
	res, err := r.DB.Exec("INSERT INTO users (name, password_hash) values (?, ?)", name, passwordHash)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
