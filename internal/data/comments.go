package data

import (
	"fmt"

	"github.com/davidpugg/stacky/internal/utils"
)

type DBComment struct {
	ID           int    `json:"id" db:"id"`
	UserID       int    `json:"user_id" db:"user_id"`
	PostID       int    `json:"post_id" db:"post_id"`
	Body         string `json:"body" db:"body"`
	UserAvatar   string `json:"user_avatar" db:"user_avatar"`
	UserUsername string `json:"user_username" db:"user_username"`
	UserEmail    string `json:"user_email" db:"user_email"`
	UserCreated  string `json:"user_created" db:"user_created"`
	CreatedAt    string `json:"created_at" db:"created_at"`
}

type Comment struct {
	ID        int     `json:"id"`
	PostID    int     `json:"post_id"`
	Body      string  `json:"body"`
	User      *DBUser `json:"user"`
	CreatedAt string  `json:"created_at"`
	IsAuthor  bool    `json:"is_author"`
}

func (d *Data) GetCommentByID(userID, commentID int) (*Comment, error) {
	query := `
		SELECT c.id, c.user_id, c.post_id, c.body, c.created_at,
		u.avatar AS user_avatar, u.username AS user_username, u.email AS user_email, u.created_at AS user_created
		FROM comments AS c
		LEFT JOIN users AS u ON u.id = c.user_id
		WHERE c.id = $1
	`

	comment := &DBComment{}
	err := d.DB.Get(comment, query, commentID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return createCommentFromDB(comment, userID), nil
}

func (d *Data) GetPostComments(userID, postID int) ([]*Comment, error) {
	query := `
		SELECT c.id, c.user_id, c.post_id, c.body, c.created_at,
		u.avatar AS user_avatar, u.username AS user_username, u.email AS user_email, u.created_at AS user_created
		FROM comments AS c
		LEFT JOIN users AS u ON u.id = c.user_id
		WHERE c.post_id = $1
	`

	comments := []*DBComment{}
	err := d.DB.Select(&comments, query, postID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	commentsToReturn := []*Comment{}
	for _, comment := range comments {
		commentsToReturn = append(commentsToReturn, createCommentFromDB(comment, userID))
	}

	return commentsToReturn, nil
}

func (d *Data) CreateComment(userID, postID int, body string) (int, error) {
	var id int
	err := d.DB.QueryRow("INSERT INTO comments(user_id, post_id, body) VALUES($1, $2, $3) RETURNING id", userID, postID, body).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return id, nil
}

func (d *Data) DeleteComment(commentID int) error {
	_, err := d.DB.Exec("DELETE FROM comments WHERE id = $1", commentID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func createCommentFromDB(commentDB *DBComment, userID int) *Comment {
	return &Comment{
		ID:     commentDB.ID,
		PostID: commentDB.PostID,
		Body:   commentDB.Body,
		User: &DBUser{
			ID:        commentDB.UserID,
			Avatar:    commentDB.UserAvatar,
			Username:  commentDB.UserUsername,
			Email:     commentDB.UserEmail,
			CreatedAt: commentDB.UserCreated,
		},
		CreatedAt: utils.FormatCreateTime(commentDB.CreatedAt),
		IsAuthor:  commentDB.UserID == userID,
	}
}
