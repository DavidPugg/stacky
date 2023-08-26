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

func (h *Handlers) registerUserRoutes(c *fiber.App) {
	r := c.Group("/users")

	r.Post("/:id/follow", middleware.Authenticate, h.followUser)
	r.Delete("/:id/follow", middleware.Authenticate, h.unfollowUser)

	r.Put("/:id", middleware.Authenticate, h.updateUser)
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

func (h *Handlers) updateUser(c *fiber.Ctx) error {
	var (
		userID      = c.Params("id")
		cropString  = c.FormValue("crop-data")
		avatar, err = c.FormFile("avatar")
		authUser    = c.Locals("AuthUser").(*middleware.UserTokenData)
		cropData    data.CropData
	)

	if err != nil {
		return utils.SendAlert(c, 400, "Invalid avatar")
	}

	if userID != strconv.Itoa(authUser.ID) {
		return utils.SendAlert(c, 403, "You cannot update another user's profile")
	}

	err = json.Unmarshal([]byte(cropString), &cropData)
	if err != nil {
		fmt.Println(err)
		return utils.SendAlert(c, 400, "Invalid crop data")
	}

	avatarID, err := h.data.SaveMediaLocally(avatar, cropData)
	if err != nil {
		fmt.Println(err)
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	avatarPath := fmt.Sprintf("%s/%s", h.mediaEndpoint, avatarID)

	if err := h.data.UpdateUser(authUser.ID, avatarPath); err != nil {
		h.data.DeleteMediaLocally(avatarID)
		return utils.SendAlert(c, 500, "Internal Server Error")
	}

	if authUser.Avatar != "" {
		avatarArr := strings.Split(authUser.Avatar, "/")
		err = h.data.DeleteMediaLocally(avatarArr[len(avatarArr)-1])
	}

	newAuthData := middleware.NewUserTokenData(
		authUser.ID,
		avatarPath,
		authUser.Username,
		authUser.Email,
	)

	session, err := h.session.Get(c)
	if err != nil {
		return utils.SendAlert(c, 500, "Error getting session")
	}

	session.Set("avatar", avatarPath)

	if err := session.Save(); err != nil {
		return utils.SendAlert(c, 500, "Error saving session")
	}

	utils.SetRedirect(c, fmt.Sprintf("/u/%s", authUser.Username))
	utils.SetAlert(c, 200, "Profile updated")
	return utils.RenderPartial(c, "navbar", newAuthData)
}
