package data

type PostLike_DB struct {
	ID        int    `json:"id" db:"id"`
	UserID    int    `json:"user_id" db:"user_id"`
	PostID    int    `json:"post_id" db:"post_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

func (d *Data) CreatePostLike(userID, postID int) error {
	_, err := d.DB.Exec("INSERT INTO post_likes (user_id, post_id) VALUES (?, ?)", userID, postID)
	if err != nil {
		return err
	}

	return nil
}

func (d *Data) DeletePostLike(userID, postID int) error {
	_, err := d.DB.Exec("DELETE FROM post_likes WHERE post_id = ?", postID)
	if err != nil {
		return err
	}

	return nil
}
