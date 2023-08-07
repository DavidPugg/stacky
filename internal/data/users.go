package data

import (
	"fmt"
	"strings"
)

type User_DB struct {
	ID        int    `json:"id" db:"id"`
	Avatar    string `json:"avatar" db:"avatar"`
	Username  string `json:"username" db:"username"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"-" db:"password"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"-" db:"updated_at"`
	Followed  bool   `json:"followed" db:"followed"`
}

func (d *Data) CreateUser(avatar, username, email, password string) (int, error) {
	result, err := d.DB.Exec("INSERT INTO users (avatar, username, email, password) VALUES (?, ?, ?, ?)", avatar, username, email, password)
	if err != nil {
		if strings.Contains(err.Error(), "users.email") {
			return 0, fmt.Errorf("users.email")
		}

		if strings.Contains(err.Error(), "users.username") {
			return 0, fmt.Errorf("users.username")
		}

		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (d *Data) GetUserByEmail(userID int, email string) (*User_DB, error) {
	query := `
	SELECT *,
	EXISTS (
		SELECT 1
		FROM follows AS f
		WHERE f.followee_id = users.id
		AND f.follower_id = ?
	) AS followed
	FROM users
	WHERE email = ?`

	user := &User_DB{}
	err := d.DB.Get(user, query, userID, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *Data) GetUserByUsername(userID int, username string) (*User_DB, error) {
	query := `
	SELECT *,
	EXISTS (
		SELECT 1
		FROM follows AS f
		WHERE f.followee_id = users.id
		AND f.follower_id = ?
	) AS followed
	FROM users
	WHERE username = ?`

	user := &User_DB{}
	err := d.DB.Get(user, query, userID, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}
