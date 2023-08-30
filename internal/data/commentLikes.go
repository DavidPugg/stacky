package data

import "fmt"

func (d *Data) CreateCommentLike(authUserID, commentID int) error {
	_, err := d.DB.Exec("INSERT INTO comment_likes (user_id, comment_id) VALUES ($1, $2)", authUserID, commentID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (d *Data) DeleteCommentLike(authUserID, commentID int) error {
	_, err := d.DB.Exec("DELETE FROM comment_likes WHERE user_id = $1 AND comment_id = $2", authUserID, commentID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
