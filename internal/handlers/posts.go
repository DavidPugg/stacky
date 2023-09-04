package handlers

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/davidpugg/stacky/internal/data"
	"github.com/davidpugg/stacky/internal/middleware"
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
)

const postsPageLimit = 5
const smallPostsPageLimit = 10
const commentsPageLimit = 15

func (h *Handlers) registerPostRoutes(c *fiber.App) {
	r := c.Group("/posts")

	r.Get("/", h.getPosts)
	r.Get("/:username", h.getUserPosts)

	r.Post("/create", middleware.Authenticate, h.createPost)
	r.Delete("/:id", middleware.Authenticate, h.deletePost)

	r.Post("/:id/like", middleware.Authenticate, h.likePost)
	r.Delete("/:id/like", middleware.Authenticate, h.unlikePost)

	r.Get("/:id/comment", h.getComments)
	r.Post("/:id/comment/:commentID", middleware.Authenticate, h.createComment)
	r.Delete("/:id/comment/:commentID", middleware.Authenticate, h.deleteComment)

	r.Get("/:id/comment/:commentID/replies", h.getCommentReplies)

	r.Post("/:id/comment/:commentID/like", middleware.Authenticate, h.likeComment)
	r.Delete("/:id/comment/:commentID/like", middleware.Authenticate, h.unlikeComment)
}

func (h *Handlers) likePost(c *fiber.Ctx) error {
	var (
		pID        = c.Params("id")
		count      = c.Query("count")
		authUserID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	likeCount, err := strconv.Atoi(count)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid count")
	}

	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	err = h.data.CreatePostLike(authUserID, postID)
	if err != nil {
		fmt.Println(err)
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	return utils.RenderPartial(c, "likeButton", fiber.Map{
		"ID":        postID,
		"Liked":     true,
		"LikeCount": likeCount + 1,
	})
}

func (h *Handlers) unlikePost(c *fiber.Ctx) error {
	var (
		pID        = c.Params("id")
		count      = c.Query("count")
		authUserID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	likeCount, err := strconv.Atoi(count)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid count")
	}

	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	err = h.data.DeletePostLike(authUserID, postID)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	return utils.RenderPartial(c, "likeButton", fiber.Map{
		"ID":        postID,
		"Liked":     false,
		"LikeCount": likeCount - 1,
	})
}

func (h *Handlers) createComment(c *fiber.Ctx) error {
	var (
		cID  = c.Params("commentID")
		pID  = c.Params("id")
		body = c.FormValue("comment")
		user = c.Locals("AuthUser").(*middleware.UserTokenData)
	)

	if len(body) == 0 {
		return utils.SendAlert(c, 400, "Invalid comment")
	}

	commentID, err := strconv.Atoi(cID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid comment ID")
	}

	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	id, err := h.data.CreateComment(user.ID, postID, commentID, body)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	comment := utils.RevealeObject[*data.Comment]{
		Data: &data.Comment{
			ID:        int(id),
			PostID:    postID,
			CommentID: commentID,
			User: &data.User{
				ID:       user.ID,
				Username: user.Username,
				Avatar:   user.Avatar,
				Email:    user.Email,
			},
			Body:      body,
			CreatedAt: "Now",
			IsAuthor:  true,
		},
		IsLast:   true,
		Page:     1,
		LastPage: true,
	}

	utils.SetTrigger(c, "swap", utils.Trigger{
		Name: "removeNoComments",
	})

	return utils.RenderPartial(c, "comment", comment)
}

func (h *Handlers) deleteComment(c *fiber.Ctx) error {
	var (
		pID        = c.Params("commentID")
		authUserID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	comment, err := h.data.GetCommentByID(authUserID, postID)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	if comment.User.ID != authUserID {
		return utils.SendAlert(c, 403, "Forbidden")
	}

	err = h.data.DeleteComment(postID)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	utils.SetTrigger(c, "swap", utils.Trigger{
		Name: "addNoComments",
	})

	return c.SendStatus(200)
}

func (h *Handlers) createPost(c *fiber.Ctx) error {
	var (
		description = c.FormValue("description")
		user        = c.Locals("AuthUser").(*middleware.UserTokenData)
		img, err    = c.FormFile("image")
		cropString  = c.FormValue("crop-data")
		cropData    utils.CropData
	)

	if err != nil {
		return utils.SendAlert(c, 400, "Invalid image")
	}

	err = json.Unmarshal([]byte(cropString), &cropData)
	if err != nil {
		fmt.Println(err)
		return utils.SendAlert(c, 400, "Invalid crop data")
	}

	ci, err := utils.CropImage(img, cropData)

	imageID, err := h.data.SaveMedia(ci, filepath.Ext(img.Filename))
	if err != nil {
		fmt.Println(err)
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	err = h.data.CreatePost(user.ID, imageID, description)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	utils.SetRedirect(c, fmt.Sprintf("/u/%s", user.Username))

	return utils.SendAlert(c, 200, "Post created")
}

func (h *Handlers) getPosts(c *fiber.Ctx) error {
	var (
		p          = c.Query("page")
		location   = strings.Split(c.Get("Referer"), "/")[3]
		authUserID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	page, err := strconv.Atoi(p)
	if err != nil {
		page = 1
	}

	posts := []*data.Post{}
	if location == "" {
		posts, err = h.data.GetFollowedPosts(authUserID, page, postsPageLimit)
	} else {
		posts, err = h.data.GetAllPosts(authUserID, page, postsPageLimit)
	}

	rp := utils.CreateRevealeObjects(posts, page, postsPageLimit)

	if err != nil {
		utils.SendAlert(c, 500, "Internal Server Error")
		return err
	}

	return utils.RenderPartial(c, "postList", rp)
}

func (h *Handlers) getUserPosts(c *fiber.Ctx) error {
	var (
		p          = c.Query("page")
		username   = c.Params("username")
		authUserID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	page, err := strconv.Atoi(p)
	if err != nil {
		page = 1
	}

	posts, err := h.data.GetPostsOfUserByUsername(authUserID, username, page, smallPostsPageLimit)
	if err != nil {
		utils.SendAlert(c, 500, "Internal Server Error")
		return err
	}

	rp := utils.CreateRevealeObjects(posts, page, postsPageLimit)

	return utils.RenderPartial(c, "smallPostList", rp)
}

func (h *Handlers) deletePost(c *fiber.Ctx) error {
	var (
		pID        = c.Params("id")
		authUserID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	post, err := h.data.GetPostByID(authUserID, postID)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	if post.User.ID != authUserID {
		return utils.SendAlert(c, 403, "Forbidden")
	}

	err = h.data.DeletePostByID(postID)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	imageArray := strings.Split(post.Image, "/")
	err = h.data.DeleteMedia(imageArray[len(imageArray)-1])
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	utils.SetRedirect(c, fmt.Sprintf("/u/%s", post.User.Username))

	return utils.SendAlert(c, 200, "Post deleted")
}

func (h *Handlers) likeComment(c *fiber.Ctx) error {
	var (
		cID        = c.Params("commentID")
		count      = c.Query("count")
		authUserID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	likeCount, err := strconv.Atoi(count)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid count")
	}

	commentID, err := strconv.Atoi(cID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid comment ID")
	}

	err = h.data.CreateCommentLike(authUserID, commentID)
	if err != nil {
		fmt.Println(err)
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	return utils.RenderPartial(c, "commentLikeButton", fiber.Map{
		"ID":        commentID,
		"PostID":    c.Params("id"),
		"Liked":     true,
		"LikeCount": likeCount + 1,
	})
}

func (h *Handlers) unlikeComment(c *fiber.Ctx) error {
	var (
		cID        = c.Params("commentID")
		count      = c.Query("count")
		authUserID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	likeCount, err := strconv.Atoi(count)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid count")
	}

	commentID, err := strconv.Atoi(cID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid comment ID")
	}

	err = h.data.DeleteCommentLike(authUserID, commentID)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	return utils.RenderPartial(c, "commentLikeButton", fiber.Map{
		"ID":        commentID,
		"PostID":    c.Params("id"),
		"Liked":     false,
		"LikeCount": likeCount - 1,
	})
}

func (h *Handlers) getCommentReplies(c *fiber.Ctx) error {
	var (
		p          = c.Query("page")
		cID        = c.Params("commentID")
		authUserID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	page, err := strconv.Atoi(p)
	if err != nil {
		page = 1
	}

	commentID, err := strconv.Atoi(cID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid comment ID")
	}

	comments, err := h.data.GetCommentReplies(authUserID, commentID, page, commentsPageLimit)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}
	rc := utils.CreateRevealeObjects(comments, page, commentsPageLimit)

	return utils.RenderPartial(c, "replies", rc)
}

func (h *Handlers) getComments(c *fiber.Ctx) error {
	var (
		p          = c.Query("page")
		pID        = c.Params("id")
		authUserID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	page, err := strconv.Atoi(p)
	if err != nil {
		page = 1
	}

	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid comment ID")
	}

	comments, err := h.data.GetPostComments(authUserID, postID, page, commentsPageLimit)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}
	rc := utils.CreateRevealeObjects(comments, page, commentsPageLimit)

	return utils.RenderPartial(c, "replies", rc)
}
