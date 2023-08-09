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
	IsAuthor  bool     `json:"is_author"`
}

const baseCommentsQuery = `
	SELECT c.id, c.user_id, c.post_id, c.body, c.created_at,
	u.avatar AS user_avatar, u.username AS user_username, u.email AS user_email, u.created_at AS user_created
	FROM comments AS c
	LEFT JOIN users AS u ON u.id = c.user_id
`

func (d *Data) GetCommentByID(userID, commentID int) (*Comment, error) {
	query := baseCommentsQuery + `WHERE c.id = ?`

	comment := &Comment_DB{}
	err := d.DB.Get(comment, query, commentID)
	if err != nil {
		return nil, err
	}

	return createCommentFromDB(comment, userID), nil
}

func (d *Data) GetPostComments(userID, postID int) ([]*Comment, error) {
	query := baseCommentsQuery + `WHERE c.post_id = ?`

	comments := []*Comment_DB{}
	err := d.DB.Select(&comments, query, postID)
	if err != nil {
		return nil, err
	}

	commentsToReturn := []*Comment{}
	for _, comment := range comments {
		commentsToReturn = append(commentsToReturn, createCommentFromDB(comment, userID))
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

	c, err := d.GetCommentByID(userID, int(commentID))
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

func createCommentFromDB(commentDB *Comment_DB, userID int) *Comment {
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
		IsAuthor:  commentDB.UserID == userID,
	}
}
