package data

import "fmt"

type Post_DB struct {
	ID           int    `json:"id" db:"id"`
	UserID       int    `json:"user_id" db:"user_id"`
	Image        string `json:"image" db:"image"`
	Description  string `json:"description" db:"description"`
	CreatedAt    string `json:"created_at" db:"created_at"`
	UserAvatar   string `json:"user_avatar" db:"user_avatar"`
	UserUsername string `json:"user_username" db:"user_username"`
	UserEmail    string `json:"user_email" db:"user_email"`
	UserCreated  string `json:"user_created" db:"user_created"`
	LikeCount    int    `json:"like_count" db:"like_count"`
}

type Post struct {
	ID          int      `json:"id" db:"id"`
	Image       string   `json:"image" db:"image"`
	Description string   `json:"description" db:"description"`
	CreatedAt   string   `json:"created_at" db:"created_at"`
	User        *User_DB `json:"user" db:"user"`
	LikeCount   int      `json:"like_count" db:"like_count"`
	// CommentCount int      `json:"comment_count" db:"comment_count"`
}

func (d *Data) GetPosts() ([]*Post, error) {
	query := `
		SELECT p.id, p.user_id, p.image, p.description, p.created_at,
		u.avatar AS user_avatar, u.username AS user_username, u.email AS user_email, u.created_at AS user_created,
		COUNT(pl.id) AS like_count
		FROM posts AS p
		LEFT JOIN users AS u ON p.user_id = u.id
		LEFT JOIN post_likes AS pl ON p.id = pl.post_id
		GROUP BY p.id
		ORDER BY created_at`

	posts := []*Post_DB{}
	err := d.DB.Select(&posts, query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	postsData := []*Post{}
	for _, p := range posts {
		postsData = append(postsData, createPostFromDB(p))
	}

	return postsData, nil
}

func (d *Data) GetPostByID(id string) (*Post, error) {
	query := `
		SELECT p.id, p.user_id, p.image, p.description, p.created_at,
		u.avatar AS user_avatar, u.username AS user_username, u.email AS user_email, u.created_at AS user_created,
		COUNT(pl.id) AS like_count
		FROM posts AS p
		LEFT JOIN users AS u ON p.user_id = u.id
		LEFT JOIN post_likes AS pl ON p.id = pl.post_id
		WHERE p.id = ?
		GROUP BY p.id`

	post := &Post_DB{}
	err := d.DB.Get(post, query, id)
	if err != nil {
		return nil, err
	}

	p := createPostFromDB(post)

	return p, nil
}

func (d *Data) CreatePost(userID int, image, description string) error {
	_, err := d.DB.Exec("INSERT INTO posts (user_id, image, description) VALUES (?, ?, ?)", userID, image, description)
	if err != nil {
		return err
	}

	return nil
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

func createPostFromDB(post *Post_DB) *Post {
	return &Post{
		ID:          post.ID,
		Image:       post.Image,
		Description: post.Description,
		CreatedAt:   post.CreatedAt,
		User: &User_DB{
			ID:        post.UserID,
			Avatar:    post.UserAvatar,
			Username:  post.UserUsername,
			Email:     post.UserEmail,
			CreatedAt: post.UserCreated,
		},
		LikeCount: post.LikeCount,
		// CommentCount: post.CommentCount,
	}
}
