package data

import (
	"fmt"
)

type Comment struct {
	ID        int    `json:"id" db:"id"`
	PostID    int    `json:"post_id" db:"post_id"`
	Body      string `json:"body" db:"body"`
	User      *User  `json:"user" db:"user"`
	CreatedAt string `json:"created_at" db:"created_at"`
	IsAuthor  bool   `json:"is_author" db:"is_author"`
}

func (d *Data) GetCommentByID(authUserID, commentID int) (*Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.body, c.created_at,
		u.id, u.avatar, u.username, u.email, u.created_at
		FROM comments AS c
		LEFT JOIN users AS u ON u.id = c.user_id
		WHERE c.id = $1
	`

	row := d.DB.QueryRow(query, commentID)

	comment, err := scanComment(row, authUserID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return comment, nil
}

func (d *Data) GetPostComments(authUserID, postID int) ([]*Comment, error) {
	var comments []*Comment

	query := `
		SELECT c.id, c.post_id, c.body, c.created_at,
		u.id, u.avatar, u.username, u.email, u.created_at
		FROM comments AS c
		LEFT JOIN users AS u ON u.id = c.user_id
		WHERE c.post_id = $1
	`

	rows, err := d.DB.Query(query, postID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		comment, err := scanComment(rows, authUserID)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (d *Data) CreateComment(authUserID, postID int, body string) (int, error) {
	var id int

	err := d.DB.QueryRow("INSERT INTO comments(user_id, post_id, body) VALUES($1, $2, $3) RETURNING id", authUserID, postID, body).Scan(&id)
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

func scanComment(row Scanner, authUserID int) (*Comment, error) {
	var (
		comment = &Comment{}
		user    = &User{}
	)

	if err := row.Scan(
		&comment.ID, &comment.PostID, &comment.Body, &comment.CreatedAt,
		&user.ID, &user.Avatar, &user.Username, &user.Email, &user.CreatedAt,
	); err != nil {
		fmt.Println(err)
		return nil, err
	}

	comment.User = user
	comment.IsAuthor = comment.User.ID == authUserID

	return comment, nil
}
