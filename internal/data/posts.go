package data

type Post_DB struct {
	ID          int    `json:"id" db:"id"`
	Image       string `json:"image" db:"image"`
	Description string `json:"description" db:"description"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"-" db:"updated_at"`
}

type Post struct {
	ID           int      `json:"id" db:"id"`
	Image        string   `json:"image" db:"image"`
	Description  string   `json:"description" db:"description"`
	CreatedAt    string   `json:"created_at" db:"created_at"`
	User         *User_DB `json:"user" db:"user"`
	LikeCount    int      `json:"like_count" db:"like_count"`
	CommentCount int      `json:"comment_count" db:"comment_count"`
}

func (d *Data) GetPosts() ([]*Post, error) {
	// posts := []*Post_DB{}
	// err := d.DB.Select(&posts, "SELECT * FROM posts ORDER BY created_at")
	// if err != nil {
	// 	return nil, err
	// }

	// return posts, nil

	const image string = "https://dfstudio-d420.kxcdn.com/wordpress/wp-content/uploads/2019/06/digital_camera_photo-1080x675.jpg"

	posts := []*Post{
		{
			ID:          1,
			Image:       image,
			Description: "This is a test post",
			CreatedAt:   "2021-09-01 00:00:00",
			User: &User_DB{
				ID:       1,
				Avatar:   image,
				Username: "testuser",
				Email:    "test@test.com",
			},
			LikeCount:    423,
			CommentCount: 23,
		},
	}

	return posts, nil
}

func (d *Data) CreatePost(image, description string) error {
	_, err := d.DB.Exec("INSERT INTO posts (image, description) VALUES (?, ?)", image, description)
	if err != nil {
		return err
	}

	return nil
}

func (d *Data) GetPostByID(id string) (*Post_DB, error) {
	post := &Post_DB{}
	err := d.DB.Get(post, "SELECT id FROM posts WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (d *Data) DeletePostByID(id string) error {
	_, err := d.DB.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (d *Data) UpdatePostByID(id, image, description string) error {
	_, err := d.DB.Exec("UPDATE posts SET image = ?, description = ? WHERE id = ?", image, description, id)
	if err != nil {
		return err
	}

	return nil
}
