package handlers

import (
	"fmt"
	"strconv"

	"github.com/davidpugg/stacky/internal/data"
	"github.com/davidpugg/stacky/internal/middleware"
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerPostRoutes(c *fiber.App) {
	r := c.Group("/posts")
	r.Post("/:id/like", middleware.Authenticate, h.likePost)
	r.Delete("/:id/like", middleware.Authenticate, h.unlikePost)
	r.Post("/:id/comment", middleware.Authenticate, h.createComment)
	r.Delete("/:id/comment", middleware.Authenticate, h.deleteComment)
	r.Post("create", middleware.Authenticate, h.createPost)
}

func (h *Handlers) likePost(c *fiber.Ctx) error {
	pID := c.Params("id")
	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	userID := c.Locals("AuthUser").(*middleware.UserTokenData).ID

	err = h.data.CreatePostLike(userID, postID)
	if err != nil {
		fmt.Println(err)
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	utils.SetTrigger(c, utils.Trigger{
		Name: "updateLikeCount",
		Data: &fiber.Map{
			"postID": postID,
			"method": "like",
		},
	})

	return utils.RenderPartial(c, "likeButton", fiber.Map{
		"ID":    postID,
		"Liked": true,
	})
}

func (h *Handlers) unlikePost(c *fiber.Ctx) error {
	pID := c.Params("id")
	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	userID := c.Locals("AuthUser").(*middleware.UserTokenData).ID

	err = h.data.DeletePostLike(userID, postID)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	utils.SetTrigger(c, utils.Trigger{
		Name: "updateLikeCount",
		Data: &fiber.Map{
			"postID": postID,
			"method": "unlike",
		},
	})

	return utils.RenderPartial(c, "likeButton", fiber.Map{
		"ID":    postID,
		"Liked": false,
	})
}

func (h *Handlers) createComment(c *fiber.Ctx) error {
	body := c.FormValue("comment")
	if len(body) == 0 {
		return utils.SendAlert(c, 400, "Invalid comment")
	}

	pID := c.Params("id")
	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	user := c.Locals("AuthUser").(*middleware.UserTokenData)

	commentID, err := h.data.CreateComment(user.ID, postID, body)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	comment := data.Comment{
		ID: int(commentID),
		User: &data.User_DB{
			ID:       user.ID,
			Username: user.Username,
			Avatar:   user.Avatar,
			Email:    user.Email,
		},
		Body:      body,
		CreatedAt: "Now",
		IsAuthor:  true,
	}

	utils.SetTrigger(c, utils.Trigger{
		Name: "removeNoComments",
	})

	return utils.RenderPartial(c, "comment", comment)
}

func (h *Handlers) deleteComment(c *fiber.Ctx) error {
	pID := c.Params("id")
	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	userID := c.Locals("AuthUser").(*middleware.UserTokenData).ID

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

	utils.SetTrigger(c, utils.Trigger{
		Name: "addNoComments",
	})

	return c.SendStatus(200)
}

func (h *Handlers) createPost(c *fiber.Ctx) error {
	description := c.FormValue("description")

	image, err := c.FormFile("image")
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid image")
	}

	path, err := h.data.CreateMediaLocally(c, image)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	user := c.Locals("AuthUser").(*middleware.UserTokenData)

	err = h.data.CreatePost(user.ID, path, description)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	utils.SetRedirect(c, fmt.Sprintf("/u/%s", user.Username))

	return utils.SendAlert(c, 200, "Post created")
}
