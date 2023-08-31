package data

import (
	"database/sql"
	"fmt"

	"github.com/davidpugg/stacky/internal/utils"
)

type Comment struct {
	ID         int    `json:"id" db:"id"`
	PostID     int    `json:"post_id" db:"post_id"`
	CommentID  int    `json:"comment_id" db:"comment_id"`
	Body       string `json:"body" db:"body"`
	User       *User  `json:"user" db:"user"`
	CreatedAt  string `json:"created_at" db:"created_at"`
	IsAuthor   bool   `json:"is_author" db:"is_author"`
	LikeCount  int    `json:"like_count" db:"like_count"`
	Liked      bool   `json:"liked" db:"liked"`
	HasReplies bool   `json:"has_replies" db:"has_replies"`
}

func createCommentQuery(q string) string {
	return fmt.Sprintf(
		`
			SELECT c.id, c.post_id, c.comment_id, c.body, c.created_at,
			u.id, u.avatar, u.username, u.email, u.created_at,
			COUNT(cl.id) AS like_count,
			EXISTS (
				SELECT 1
				FROM comment_likes AS cl2
				WHERE cl2.comment_id = c.id
				AND cl2.user_id = $1
			) AS liked,
			EXISTS (
				SELECT 1
				FROM comments AS c2
				WHERE c2.comment_id = c.id
			) AS has_replies
			FROM comments AS c
			LEFT JOIN users AS u ON u.id = c.user_id
			LEFT JOIN comment_likes AS cl ON cl.comment_id = c.id
			%s
			GROUP BY c.id, u.id
			ORDER BY c.created_at DESC
	`,
		q,
	)
}

func (d *Data) GetCommentByID(authUserID, commentID int) (*Comment, error) {
	query := createCommentQuery("WHERE c.id = $2")

	row := d.DB.QueryRow(query, authUserID, commentID)

	comment, err := scanComment(row, authUserID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return comment, nil
}

func (d *Data) GetPostComments(authUserID, postID int) ([]*Comment, error) {
	var comments []*Comment

	query := createCommentQuery("WHERE c.post_id = $2 AND c.comment_id IS NULL")

	rows, err := d.DB.Query(query, authUserID, postID)
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

func (d *Data) GetCommentReplies(authUserID, commentID int) ([]*Comment, error) {
	var comments []*Comment

	query := createCommentQuery("WHERE c.comment_id = $2")

	rows, err := d.DB.Query(query, authUserID, commentID)
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

func (d *Data) CreateComment(authUserID, postID, commentID int, body string) (int, error) {
	var id int

	var cID = interface{}(nil)
	if commentID != 0 {
		cID = commentID
	}

	err := d.DB.QueryRow("INSERT INTO comments(user_id, post_id, comment_id, body) VALUES($1, $2, $3, $4) RETURNING id", authUserID, postID, cID, body).Scan(&id)
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
		commentID = sql.NullInt64{}
		comment   = &Comment{}
		user      = &User{}
	)

	if err := row.Scan(
		&comment.ID, &comment.PostID, &commentID, &comment.Body, &comment.CreatedAt,
		&user.ID, &user.Avatar, &user.Username, &user.Email, &user.CreatedAt,
		&comment.LikeCount, &comment.Liked, &comment.HasReplies,
	); err != nil {
		fmt.Println(err)
		return nil, err
	}

	if commentID.Valid {
		comment.CommentID = int(commentID.Int64)
	}

	comment.User = user
	comment.User.Avatar = utils.CreateImagePath(comment.User.Avatar)
	comment.CreatedAt = utils.FormatCreateTime(comment.CreatedAt)
	comment.IsAuthor = comment.User.ID == authUserID

	return comment, nil
}
