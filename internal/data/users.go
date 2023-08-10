package data

import (
	"fmt"
	"strings"
)

type User_DB struct {
	ID             int    `json:"id" db:"id"`
	Avatar         string `json:"avatar" db:"avatar"`
	Username       string `json:"username" db:"username"`
	Email          string `json:"email" db:"email"`
	Password       string `json:"-" db:"password"`
	CreatedAt      string `json:"created_at" db:"created_at"`
	UpdatedAt      string `json:"-" db:"updated_at"`
	Followed       bool   `json:"followed" db:"followed"`
	FollowingCount int    `json:"following_count" db:"following_count"`
	FollowersCount int    `json:"followers_count" db:"followers_count"`
}

type UserWithPosts struct {
	User_DB
	Posts []*Post `json:"posts"`
}

const baseUsersQuery = `
	SELECT users.id, users.avatar, users.username, users.password, users.email, users.created_at, users.updated_at,
	EXISTS (
		SELECT 1
		FROM follows AS f
		WHERE f.followee_id = users.id
		AND f.follower_id = ?
	) AS followed `

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

func (d *Data) GetUserByID(authUserID, userID int) (*User_DB, error) {
	query := baseUsersQuery + `FROM users WHERE id = ?`

	user := &User_DB{}
	err := d.DB.Get(user, query, authUserID, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *Data) GetUserByEmail(userID int, email string) (*User_DB, error) {
	query := baseUsersQuery + `FROM users WHERE email = ?`

	user := &User_DB{}
	err := d.DB.Get(user, query, userID, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *Data) GetUserByUsername(userID int, username string) (*User_DB, error) {
	query := baseUsersQuery + `FROM users WHERE username = ?`

	user := &User_DB{}
	err := d.DB.Get(user, query, userID, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *Data) GetUserWithPostsByUsername(userID int, username string) (*UserWithPosts, error) {
	postsChan := make(chan []*Post)

	go func() {
		posts, err := d.GetPostsOfUserByUsername(userID, username)
		if err != nil {
			postsChan <- nil
			return
		}

		postsChan <- posts
	}()

	query := baseUsersQuery + `,
	COUNT(f.id) AS followers_count,
	COUNT(f2.id) AS following_count
	FROM users
	LEFT JOIN follows AS f ON f.followee_id = users.id
	LEFT JOIN follows AS f2 ON f2.follower_id = users.id
	WHERE username = ?
	GROUP BY users.id
	`

	user := &User_DB{}
	err := d.DB.Get(user, query, userID, username)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	posts := <-postsChan
	if posts == nil {
		return nil, err
	}

	return &UserWithPosts{
		User_DB: *user,
		Posts:   posts,
	}, nil
}
