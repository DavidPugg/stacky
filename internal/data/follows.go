package data

import "fmt"

func (d *Data) CreateFollow(followerID, followeeID int) error {
	_, err := d.DB.Exec("INSERT INTO follows (follower_id, followee_id) VALUES ($1, $2)", followerID, followeeID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (d *Data) DeleteFollow(followerID, followeeID int) error {
	_, err := d.DB.Exec("DELETE FROM follows WHERE follower_id = $1 AND followee_id = $2", followerID, followeeID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
