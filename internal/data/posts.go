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
	Liked        bool   `json:"liked" db:"liked"`
	CommentCount int    `json:"comment_count" db:"comment_count"`
	Followed     bool   `json:"followed" db:"followed"`
}

type Post struct {
	ID           int      `json:"id" db:"id"`
	Image        string   `json:"image" db:"image"`
	Description  string   `json:"description" db:"description"`
	CreatedAt    string   `json:"created_at" db:"created_at"`
	User         *User_DB `json:"user" db:"user"`
	LikeCount    int      `json:"like_count" db:"like_count"`
	Liked        bool     `json:"liked" db:"liked"`
	CommentCount int      `json:"comment_count" db:"comment_count"`
}

type PostWithComments struct {
	Post
	Comments []*Comment `json:"comments"`
}

func (d *Data) GetPostsOfUserByUsername(userID int, username string) ([]*Post, error) {
	query := `
		SELECT p.id, p.user_id, p.image, p.description, p.created_at,
		u.avatar AS user_avatar, u.username AS user_username, u.email AS user_email, u.created_at AS user_created,
		(SELECT COUNT(*) FROM post_likes WHERE post_id = p.id) AS like_count,
		COUNT(c.id) AS comment_count,
		EXISTS (
			SELECT 1
			FROM post_likes AS pl2
			WHERE pl2.post_id = p.id
			AND pl2.user_id = ?
		) AS liked,
		EXISTS (
			SELECT 1
			FROM follows AS f
			WHERE f.followee_id = p.user_id
			AND f.follower_id = ?
		) AS followed
		FROM posts AS p
		LEFT JOIN users AS u ON p.user_id = u.id
		LEFT JOIN comments AS c ON p.id = c.post_id
		WHERE u.username = ?
		GROUP BY p.id
		ORDER BY created_at`

	posts := []*Post_DB{}
	err := d.DB.Select(&posts, query, userID, userID, username)
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

func (d *Data) GetFollowedPosts(userID int) ([]*Post, error) {
	query := `
		SELECT p.id, p.user_id, p.image, p.description, p.created_at,
		u.avatar AS user_avatar, u.username AS user_username, u.email AS user_email, u.created_at AS user_created,
		(SELECT COUNT(*) FROM post_likes WHERE post_id = p.id) AS like_count,
		COUNT(c.id) AS comment_count,
		EXISTS (
			SELECT 1
			FROM post_likes AS pl2
			WHERE pl2.post_id = p.id
			AND pl2.user_id = ?
		) AS liked,
		EXISTS (
			SELECT 1
			FROM follows AS f
			WHERE f.followee_id = p.user_id
			AND f.follower_id = ?
		) AS followed
		FROM posts AS p
		LEFT JOIN users AS u ON p.user_id = u.id
		LEFT JOIN comments AS c ON p.id = c.post_id
		LEFT JOIN follows AS f ON p.user_id = f.followee_id
		WHERE f.follower_id = ?
		GROUP BY p.id
		ORDER BY created_at`

	posts := []*Post_DB{}
	err := d.DB.Select(&posts, query, userID, userID, userID)
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

func (d *Data) GetAllPosts(userID int) ([]*Post, error) {
	query := `
		SELECT p.id, p.user_id, p.image, p.description, p.created_at,
		u.avatar AS user_avatar, u.username AS user_username, u.email AS user_email, u.created_at AS user_created,
		(SELECT COUNT(*) FROM post_likes WHERE post_id = p.id) AS like_count,
		COUNT(c.id) AS comment_count,
		EXISTS (
			SELECT 1
			FROM post_likes AS pl2
			WHERE pl2.post_id = p.id
			AND pl2.user_id = ?
		) AS liked,
		EXISTS (
			SELECT 1
			FROM follows AS f
			WHERE f.followee_id = p.user_id
			AND f.follower_id = ?
		) AS followed
		FROM posts AS p
		LEFT JOIN users AS u ON p.user_id = u.id
		LEFT JOIN comments AS c ON p.id = c.post_id
		GROUP BY p.id
		ORDER BY created_at`

	posts := []*Post_DB{}
	err := d.DB.Select(&posts, query, userID, userID)
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

func (d *Data) GetPostWithCommentsByID(userID, postID int) (*PostWithComments, error) {
	query := `
		SELECT p.id, p.user_id, p.image, p.description, p.created_at,
		u.avatar AS user_avatar, u.username AS user_username, u.email AS user_email, u.created_at AS user_created,
		(SELECT COUNT(*) FROM post_likes WHERE post_id = p.id) AS like_count,
		COUNT(c.id) AS comment_count,
		EXISTS (
			SELECT 1
			FROM post_likes AS pl2
			WHERE pl2.post_id = p.id
			AND pl2.user_id = ?
		) AS liked,
		EXISTS (
			SELECT 1
			FROM follows AS f
			WHERE f.followee_id = p.user_id
			AND f.follower_id = ?
		) AS followed
		FROM posts AS p
		LEFT JOIN users AS u ON p.user_id = u.id
		LEFT JOIN comments AS c ON p.id = c.post_id
		WHERE p.id = ?
		GROUP BY p.id`

	post := &Post_DB{}
	err := d.DB.Get(post, query, userID, userID, postID)
	if err != nil {
		return nil, err
	}

	commentChan := make(chan []*Comment)

	go func() {
		comments, err := d.GetPostComments(userID, postID)
		if err != nil {
			commentChan <- nil
			return
		}

		commentChan <- comments
	}()

	p := createPostFromDB(post)

	comments := <-commentChan
	if comments == nil {
		return nil, err
	}

	pwc := createPostWithComments(p, comments)

	return pwc, nil
}

func (d *Data) CreatePost(userID int, image, description string) error {
	_, err := d.DB.Exec("INSERT INTO posts (user_id, image, description) VALUES (?, ?, ?)", userID, image, description)
	if err != nil {
		return err
	}

	return nil
}

func (d *Data) DeletePostByID(id int) error {
	_, err := d.DB.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (d *Data) UpdatePostByID(id int, image, description string) error {
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
			Followed:  post.Followed,
		},
		LikeCount:    post.LikeCount,
		Liked:        post.Liked,
		CommentCount: post.CommentCount,
	}
}

func createPostWithComments(post *Post, comments []*Comment) *PostWithComments {
	return &PostWithComments{
		Post:     *post,
		Comments: comments,
	}
}
