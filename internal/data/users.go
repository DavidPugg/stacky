package data

import (
	"fmt"
	"strings"

	"github.com/davidpugg/stacky/internal/utils"
)

type User struct {
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
	User
	Posts []*Post `json:"posts"`
}

func createUserQuery(q string) string {
	return fmt.Sprintf(
		`
			SELECT users.id, users.avatar, users.username, users.password, users.email, users.created_at, users.updated_at,
			COUNT(f.id) AS followers_count,
			COUNT(f2.id) AS following_count,
			EXISTS (
				SELECT 1
				FROM follows AS f
				WHERE f.followee_id = users.id
				AND f.follower_id = $1
			) AS followed
			FROM users
			LEFT JOIN follows AS f ON f.followee_id = users.id
			LEFT JOIN follows AS f2 ON f2.follower_id = users.id
			%s
			GROUP BY users.id
	`,
		q,
	)
}

func (d *Data) CreateUser(avatar, username, email, password string) (int, error) {
	var id int

	err := d.DB.QueryRow("INSERT INTO users (avatar, username, email, password) VALUES ($1, $2, $3, $4) RETURNING id", avatar, username, email, password).Scan(&id)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "email_unique") {
			return 0, fmt.Errorf("email_unique")
		}

		if strings.Contains(err.Error(), "username_unique") {
			return 0, fmt.Errorf("username_unique")
		}

		return 0, err
	}

	return id, nil
}

func (d *Data) GetUserByID(authUserID, userID int) (*User, error) {
	var user = &User{}

	query := createUserQuery("WHERE id = $2")

	err := d.DB.Get(user, query, authUserID, userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user.Avatar = utils.CreateImagePath(user.Avatar)

	return user, nil
}

func (d *Data) GetUserByEmail(authUserID int, email string) (*User, error) {
	var user = &User{}

	query := createUserQuery("WHERE email = $2")

	err := d.DB.Get(user, query, authUserID, email)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user.Avatar = utils.CreateImagePath(user.Avatar)

	return user, nil
}

func (d *Data) GetUserByUsername(authUserID int, username string) (*User, error) {
	var user = &User{}

	query := createUserQuery("WHERE username = $2")

	err := d.DB.Get(user, query, authUserID, username)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user.Avatar = utils.CreateImagePath(user.Avatar)

	return user, nil
}

func (d *Data) GetUserWithPostsByUsername(authUserID int, username string) (*UserWithPosts, error) {
	var (
		user      = &User{}
		postsChan = make(chan []*Post)
		errorChan = make(chan error)
	)

	go func() {
		posts, err := d.GetPostsOfUserByUsername(authUserID, username)
		if err != nil {
			errorChan <- err
			return
		}

		errorChan <- nil
		postsChan <- posts
	}()

	query := createUserQuery("WHERE username = $2")

	err := d.DB.Get(user, query, authUserID, username)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user.Avatar = utils.CreateImagePath(user.Avatar)

	err = <-errorChan
	if err != nil {
		return nil, err
	}

	posts := <-postsChan

	return &UserWithPosts{
		User:  *user,
		Posts: posts,
	}, nil
}

func (d *Data) UpdateUser(userID int, avatar string) error {
	_, err := d.DB.Exec("UPDATE users SET avatar = $1 WHERE id = $2", avatar, userID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
