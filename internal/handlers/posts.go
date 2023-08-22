package handlers

import (
	"encoding/json"
	"fmt"
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

	r.Post("/:id/like", middleware.Authenticate, h.likePost)
	r.Delete("/:id/like", middleware.Authenticate, h.unlikePost)

	r.Post("/:id/comment", middleware.Authenticate, h.createComment)
	r.Delete("/:id/comment", middleware.Authenticate, h.deleteComment)
}

func (h *Handlers) likePost(c *fiber.Ctx) error {
	var (
		pID    = c.Params("id")
		count  = c.Query("count")
		userID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	likeCount, err := strconv.Atoi(count)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid count")
	}

	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	err = h.data.CreatePostLike(userID, postID)
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
		pID    = c.Params("id")
		count  = c.Query("count")
		userID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	likeCount, err := strconv.Atoi(count)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid count")
	}

	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	err = h.data.DeletePostLike(userID, postID)
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

	commentID, err := h.data.CreateComment(user.ID, postID, body)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	comment := data.Comment{
		ID: int(commentID),
		User: &data.DBUser{
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
		pID    = c.Params("id")
		userID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	comment, err := h.data.GetCommentByID(userID, postID)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	if comment.User.ID != userID {
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
		cropData    = c.FormValue("crop-data")
		data        data.CropData
	)

	if err != nil {
		return utils.SendAlert(c, 400, "Invalid image")
	}

	err = json.Unmarshal([]byte(cropData), &data)
	if err != nil {
		fmt.Println(err)
		return utils.SendAlert(c, 400, "Invalid crop data")
	}

	id, err := h.data.SaveMediaLocally(img, data)
	if err != nil {
		fmt.Println(err)
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	err = h.data.CreatePost(user.ID, fmt.Sprintf("%s/%s", h.mediaEndpoint, id), description)
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
