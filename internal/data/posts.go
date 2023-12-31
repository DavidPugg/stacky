package data

import (
	"fmt"

	"github.com/davidpugg/stacky/internal/utils"
)

type Post struct {
	ID           int    `json:"id" db:"id"`
	Image        string `json:"image" db:"image"`
	Description  string `json:"description" db:"description"`
	CreatedAt    string `json:"created_at" db:"created_at"`
	User         *User  `json:"user" db:"user"`
	LikeCount    int    `json:"like_count" db:"like_count"`
	Liked        bool   `json:"liked" db:"liked"`
	CommentCount int    `json:"comment_count" db:"comment_count"`
	TotalCount   int    `json:"total_count" db:"total_count"`
}

func (p *Post) GetTotalCount() int {
	return p.TotalCount
}

const basePostQuery = `
	SELECT p.id, p.image, p.description, p.created_at,
	u.id, u.avatar, u.username, u.email, u.created_at,
	(SELECT COUNT(*) FROM post_likes WHERE post_id = p.id) AS like_count,
	COUNT(c.id) AS comment_count,
	COUNT(*) OVER() AS total_count,
	EXISTS (
		SELECT 1
		FROM post_likes AS pl2
		WHERE pl2.post_id = p.id
		AND pl2.user_id = $1
	) AS liked,
	EXISTS (
		SELECT 1
		FROM follows AS f
		WHERE f.followee_id = p.user_id
		AND f.follower_id = $2
	) AS followed
	FROM posts AS p
	LEFT JOIN users AS u ON p.user_id = u.id
	LEFT JOIN comments AS c ON p.id = c.post_id
`

func (d *Data) GetPostByID(userID, postID int) (*Post, error) {
	query := basePostQuery + `
		WHERE p.id = $3
		GROUP BY p.id, u.id
	`

	row := d.DB.QueryRow(query, userID, userID, postID)

	post, err := scanPost(row)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return post, nil
}

func (d *Data) GetPostsOfUserByUsername(userID int, username string, page, limit int) ([]*Post, error) {
	var posts []*Post

	query := basePostQuery + `
		WHERE u.username = $3
		GROUP BY p.id, u.id
		ORDER BY p.created_at DESC
		LIMIT $4 OFFSET $5
	`

	rows, err := d.DB.Query(query, userID, userID, username, limit, (page-1)*limit)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (d *Data) GetFollowedPosts(userID, page, limit int) ([]*Post, error) {
	var posts []*Post

	query := basePostQuery + `
		LEFT JOIN follows AS f ON p.user_id = f.followee_id
		WHERE f.follower_id = $3
		GROUP BY p.id, u.id
		ORDER BY p.created_at DESC
		LIMIT $4 OFFSET $5
	`

	rows, err := d.DB.Query(query, userID, userID, userID, limit, (page-1)*limit)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (d *Data) GetAllPosts(userID, page, limit int) ([]*Post, error) {
	var posts []*Post

	query := basePostQuery + `
		WHERE p.user_id != $3
		GROUP BY p.id, u.id
		ORDER BY p.created_at DESC
		LIMIT $4 OFFSET $5
	`

	rows, err := d.DB.Query(query, userID, userID, userID, limit, (page-1)*limit)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (d *Data) CreatePost(userID int, image, description string) error {
	_, err := d.DB.Exec("INSERT INTO posts (user_id, image, description) VALUES ($1, $2, $3)", userID, image, description)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (d *Data) DeletePostByID(id int) error {
	_, err := d.DB.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (d *Data) UpdatePostByID(id int, image, description string) error {
	_, err := d.DB.Exec("UPDATE posts SET image = $1, description = $2 WHERE id = $3", image, description, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func scanPost(row Scanner) (*Post, error) {
	var (
		post = &Post{}
		user = &User{}
	)

	if err := row.Scan(
		&post.ID,
		&post.Image,
		&post.Description,
		&post.CreatedAt,
		&user.ID,
		&user.Avatar,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&post.LikeCount,
		&post.CommentCount,
		&post.TotalCount,
		&post.Liked,
		&user.Followed,
	); err != nil {
		return nil, err
	}

	post.User = user
	post.Image = utils.CreateImagePath(post.Image)
	user.Avatar = utils.CreateImagePath(user.Avatar)

	return post, nil
}
