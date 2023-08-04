package data

type Comment_DB struct {
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
	ID        int      `json:"id"`
	PostID    int      `json:"post_id"`
	Body      string   `json:"body"`
	User      *User_DB `json:"user"`
	CreatedAt string   `json:"created_at"`
}

func (d *Data) GetCommentByID(commentID int) (*Comment, error) {
	var comment Comment_DB
	err := d.DB.QueryRow("SELECT id, user_id, post_id, body, created_at FROM comments WHERE id = ?", commentID).Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Body, &comment.CreatedAt)
	if err != nil {
		return nil, err
	}

	return createCommentFromDB(&comment), nil
}

func (d *Data) GetPostComments(postID string) ([]*Comment, error) {
	comments := []*Comment_DB{}
	err := d.DB.Select(&comments, "SELECT id, user_id, post_id, body, created_at FROM comments WHERE post_id = ?", postID)
	if err != nil {
		return nil, err
	}

	commentsToReturn := []*Comment{}
	for _, comment := range comments {
		commentsToReturn = append(commentsToReturn, createCommentFromDB(comment))
	}

	return commentsToReturn, nil
}

func (d *Data) CreateComment(userID, postID int, body string) (*Comment, error) {
	t, err := d.DB.Exec("INSERT INTO comments(user_id, post_id, body) VALUES(?, ?, ?)", userID, postID, body)
	if err != nil {
		return nil, err
	}

	commentID, err := t.LastInsertId()
	if err != nil {
		return nil, err
	}

	c, err := d.GetCommentByID(int(commentID))
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (d *Data) DeleteComment(commentID int) error {
	_, err := d.DB.Exec("DELETE FROM comments WHERE id = ?", commentID)
	if err != nil {
		return err
	}

	return nil
}

func createCommentFromDB(commentDB *Comment_DB) *Comment {
	return &Comment{
		ID:     commentDB.ID,
		PostID: commentDB.PostID,
		Body:   commentDB.Body,
		User: &User_DB{
			ID:        commentDB.UserID,
			Avatar:    commentDB.UserAvatar,
			Username:  commentDB.UserUsername,
			Email:     commentDB.UserEmail,
			CreatedAt: commentDB.UserCreated,
		},
		CreatedAt: commentDB.CreatedAt,
	}
}
