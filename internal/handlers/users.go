package handlers

import (
	"strconv"

	"github.com/davidpugg/stacky/internal/middleware"
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerUserRoutes(c *fiber.App) {
	r := c.Group("/users")
	r.Post("/:id/follow", middleware.Authenticate, h.followUser)
	r.Delete("/:id/follow", middleware.Authenticate, h.unfollowUser)
}

func (h *Handlers) followUser(c *fiber.Ctx) error {
	fID := c.Params("id")
	followeeID, err := strconv.Atoi(fID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid user ID")
	}

	followerID := c.Locals("AuthUser").(*middleware.UserTokenData).ID

	if followerID == followeeID {
		return utils.SendAlert(c, 400, "You cannot follow yourself")
	}

	err = h.data.CreateFollow(followerID, followeeID)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	utils.SetTrigger(c, utils.Trigger{
		Name: "updateFollowCount",
		Data: &fiber.Map{
			"followeeID": followeeID,
			"method":     "follow",
			"buttonText": "Following",
		},
	})

	return utils.RenderPartial(c, "followButton", fiber.Map{
		"ID":       followeeID,
		"Followed": true,
	})
}

func (h *Handlers) unfollowUser(c *fiber.Ctx) error {
	fID := c.Params("id")
	followeeID, err := strconv.Atoi(fID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid user ID")
	}

	followerID := c.Locals("AuthUser").(*middleware.UserTokenData).ID

	if followeeID == followerID {
		return utils.SendAlert(c, 400, "You can't unfollow yourself")
	}

	err = h.data.DeleteFollow(followerID, followeeID)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	utils.SetTrigger(c, utils.Trigger{
		Name: "updateFollowCount",
		Data: &fiber.Map{
			"followeeID": followeeID,
			"method":     "unfollow",
			"buttonText": "Follow",
		},
	})

	return utils.RenderPartial(c, "followButton", fiber.Map{
		"ID":       followeeID,
		"Followed": false,
	})
}
