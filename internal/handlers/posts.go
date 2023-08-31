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

func (h *Handlers) registerPostRoutes(c *fiber.App) {
	r := c.Group("/posts")

	r.Get("/", h.getPosts)
	r.Post("/create", middleware.Authenticate, h.createPost)
	r.Delete("/:id", middleware.Authenticate, h.deletePost)

	r.Post("/:id/like", middleware.Authenticate, h.likePost)
	r.Delete("/:id/like", middleware.Authenticate, h.unlikePost)

	r.Post("/:id/comment", middleware.Authenticate, h.createComment)
	r.Delete("/:id/comment/:commentID", middleware.Authenticate, h.deleteComment)

	r.Get("/:id/comment/:commentID/replies", h.getCommentReplies)
	r.Post("/:id/comment/:commentID/replies", h.createReply)

	r.Post("/:id/comment/like", middleware.Authenticate, h.likeComment)
	r.Delete("/:id/comment/like", middleware.Authenticate, h.unlikeComment)
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
		pID  = c.Params("id")
		body = c.FormValue("comment")
		user = c.Locals("AuthUser").(*middleware.UserTokenData)
	)

	if len(body) == 0 {
		return utils.SendAlert(c, 400, "Invalid comment")
	}

	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	commentID, err := h.data.CreateComment(user.ID, postID, 0, body)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	comment := data.Comment{
		ID: int(commentID),
		User: &data.User{
			ID:       user.ID,
			Username: user.Username,
			Avatar:   user.Avatar,
			Email:    user.Email,
		},
		Body:      body,
		CreatedAt: "Now",
		IsAuthor:  true,
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
		location = strings.Split(c.Get("Referer"), "/")[3]
		user     = c.Locals("AuthUser").(*middleware.UserTokenData)
	)

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}

	var posts []*data.LastPost

	if location == "discover" {
		posts, err = h.data.GetAllPosts(user.ID, page)
	} else if location == "" {
		posts, err = h.data.GetFollowedPosts(user.ID, page)
	}

	if err != nil {
		utils.SendAlert(c, 500, "Internal Server Error")
		return err
	}

	return utils.RenderPartial(c, "postList", posts)
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
		cID        = c.Params("id")
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
		"Liked":     true,
		"LikeCount": likeCount + 1,
	})
}

func (h *Handlers) unlikeComment(c *fiber.Ctx) error {
	var (
		cID        = c.Params("id")
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
		"Liked":     false,
		"LikeCount": likeCount - 1,
	})
}

func (h *Handlers) getCommentReplies(c *fiber.Ctx) error {
	var (
		cID        = c.Params("commentID")
		authUserID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	commentID, err := strconv.Atoi(cID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid comment ID")
	}

	comments, err := h.data.GetCommentReplies(authUserID, commentID)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	return utils.RenderPartial(c, "replies", comments)
}

func (h *Handlers) createReply(c *fiber.Ctx) error {
	var (
		pID  = c.Params("id")
		cID  = c.Params("commentID")
		body = c.FormValue("comment")
		user = c.Locals("AuthUser").(*middleware.UserTokenData)
	)

	if len(body) == 0 {
		return utils.SendAlert(c, 400, "Invalid comment")
	}

	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	commentID, err := strconv.Atoi(cID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid comment ID")
	}

	replyID, err := h.data.CreateComment(user.ID, postID, commentID, body)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	reply := data.Comment{
		ID: int(replyID),
		User: &data.User{
			ID:       user.ID,
			Username: user.Username,
			Avatar:   user.Avatar,
			Email:    user.Email,
		},
		Body:      body,
		CreatedAt: "Now",
		IsAuthor:  true,
	}

	return utils.RenderPartial(c, "comment", reply)
}
