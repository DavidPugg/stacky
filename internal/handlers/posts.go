package handlers

import (
	"fmt"
	"strconv"

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

	return utils.RenderPartial(c, "likeButton", fiber.Map{
		"ID":    postID,
		"Liked": false,
	})
}

func (h *Handlers) createComment(c *fiber.Ctx) error {
	commentForm := c.FormValue("comment")
	if len(commentForm) == 0 {
		return utils.SendAlert(c, 400, "Invalid comment")
	}

	pID := c.Params("id")
	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	userID := c.Locals("AuthUser").(*middleware.UserTokenData).ID

	comment, err := h.data.CreateComment(userID, postID, commentForm)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
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
