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
	var (
		userID     = c.Params("id")
		authUserID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	followeeID, err := strconv.Atoi(userID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid user ID")
	}

	if authUserID == followeeID {
		return utils.SendAlert(c, 400, "You cannot follow yourself")
	}

	err = h.data.CreateFollow(authUserID, followeeID)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	utils.SetTrigger(c, "swap", utils.Trigger{
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
	var (
		userID     = c.Params("id")
		authUserId = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	followeeID, err := strconv.Atoi(userID)
	if err != nil {
		return utils.SendAlert(c, 400, "Invalid user ID")
	}

	if followeeID == authUserId {
		return utils.SendAlert(c, 400, "You can't unfollow yourself")
	}

	err = h.data.DeleteFollow(authUserId, followeeID)
	if err != nil {
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	utils.SetTrigger(c, "swap", utils.Trigger{
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
