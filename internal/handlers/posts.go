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
}

func (h *Handlers) likePost(c *fiber.Ctx) error {
	pID := c.Params("id")
	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	userID := c.Locals("User").(*middleware.UserTokenData).ID

	err = h.data.CreatePostLike(userID, postID)
	if err != nil {
		fmt.Println(err)
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	return nil
}

func (h *Handlers) unlikePost(c *fiber.Ctx) error {
	pID := c.Params("id")
	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid post ID")
	}

	userID := c.Locals("User").(*middleware.UserTokenData).ID

	err = h.data.DeletePostLike(userID, postID)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	return nil
}
