package models

type User struct {
	Id           int64
	Name         string `json:"name"`
	Password     string `json:"password"`
	PasswordHash string
}
