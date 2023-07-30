package data

import (
	"fmt"
	"strings"
)

type User_DB struct {
	ID int `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"-" db:"updated_at"`
}

func (d *Data) CreateUser(avatar, username, email, password string) error {
	_, err := d.DB.Exec("INSERT INTO users (avatar, username, email, password) VALUES (?, ?, ?, ?)", avatar, username, email, password)
	if err != nil {
		if strings.Contains(err.Error(), "users.email") {
			return fmt.Errorf("users.email")
		}

		if strings.Contains(err.Error(), "users.username") {
			return fmt.Errorf("users.username")
		}

		return err
	}

	return nil
}