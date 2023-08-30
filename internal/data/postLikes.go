package data

import "fmt"

func (d *Data) CreatePostLike(authUserID, postID int) error {
	_, err := d.DB.Exec("INSERT INTO post_likes (user_id, post_id) VALUES ($1, $2)", authUserID, postID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (d *Data) DeletePostLike(authUserID, postID int) error {
	_, err := d.DB.Exec("DELETE FROM post_likes WHERE user_id = $1 AND post_id = $2", authUserID, postID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
