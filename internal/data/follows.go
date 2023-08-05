package data

type Follow_DB struct {
	ID         int `json:"id" db:"id"`
	FollowerID int `json:"follower_id" db:"follower_id"`
	FolloweeID int `json:"followee_id" db:"followee_id"`
}

func (d *Data) CreateFollow(followerID, followeeID int) error {
	_, err := d.DB.Exec("INSERT INTO follows (follower_id, followee_id) VALUES (?, ?)", followerID, followeeID)
	if err != nil {
		return err
	}

	return nil
}

func (d *Data) DeleteFollow(followerID, followeeID int) error {
	_, err := d.DB.Exec("DELETE FROM follows WHERE follower_id = ? AND followee_id = ?", followerID, followeeID)
	if err != nil {
		return err
	}

	return nil
}
