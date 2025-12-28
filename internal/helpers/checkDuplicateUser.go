package helpers

import (
	"database/sql"

	"main.go/internal/db"
)

func CheckDuplicateUser(name string) (bool, error) {
	var existing string
	err := db.DB.QueryRow("SELECT name FROM users WHERE name = ?", name).Scan(&existing)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
