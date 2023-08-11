package data

import "fmt"

type PostLike_DB struct {
	ID        int    `json:"id" db:"id"`
	UserID    int    `json:"user_id" db:"user_id"`
	PostID    int    `json:"post_id" db:"post_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

func (d *Data) CreatePostLike(userID, postID int) error {
	_, err := d.DB.Exec("INSERT INTO post_likes (user_id, post_id) VALUES ($1, $2)", userID, postID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (d *Data) DeletePostLike(userID, postID int) error {
	_, err := d.DB.Exec("DELETE FROM post_likes WHERE user_id = $1 AND post_id = $2", userID, postID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
