package handlers

import (
	"fmt"
	"strconv"

	"github.com/davidpugg/stacky/internal/middleware"
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerViewRoutes(c *fiber.App) {
	c.Get("/", middleware.MainAuthGuard, h.renderMain)
	c.Get("/discover", h.renderDiscover)
	c.Get("/register", h.renderRegister)
	c.Get("/login", middleware.LoginAuthGuard, h.renderLogin)
	c.Get("/post/:id", h.renderPost)
	c.Get("/u/:username", h.renderUser)
	c.Get("/create", middleware.MainAuthGuard, h.renderCreatePost)
	c.Get("/u/:username/edit", middleware.Authenticate, h.renderEditUser)
}

func (h *Handlers) renderMain(c *fiber.Ctx) error {
	var (
		userID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	posts, err := h.data.GetFollowedPosts(userID, 1)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching posts")
	}

	return utils.RenderPage(c, "index", fiber.Map{"Posts": posts}, &utils.PageDetails{
		Title:       "Stacky",
		Description: "Stacky is a simple social media platform",
	})
}

func (h *Handlers) renderDiscover(c *fiber.Ctx) error {
	var (
		userID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	posts, err := h.data.GetAllPosts(userID, 1)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching posts")
	}

	return utils.RenderPage(c, "index", fiber.Map{"Posts": posts}, &utils.PageDetails{
		Title:       "Stacky",
		Description: "Stacky is a simple social media platform",
	})
}

func (h *Handlers) renderLogin(c *fiber.Ctx) error {
	return utils.RenderPage(c, "login", fiber.Map{}, &utils.PageDetails{
		Title:       "Login",
		Description: "Login to Stacky",
	})
}

func (h *Handlers) renderRegister(c *fiber.Ctx) error {
	return utils.RenderPage(c, "register", fiber.Map{}, &utils.PageDetails{
		Title:       "Register",
		Description: "Register for Stacky",
	})
}

func (h *Handlers) renderPost(c *fiber.Ctx) error {
	var (
		pID    = c.Params("id")
		userID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	if pID == "" {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Invalid post ID")
	}

	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Invalid post ID")
	}

	post, err := h.data.GetPostWithCommentsByID(userID, postID)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching posts")
	}

	return utils.RenderPage(c, "post", fiber.Map{"Post": post, "ShowFollowButton": userID != post.User.ID}, &utils.PageDetails{
		Title:       fmt.Sprintf("%s - %d - Stacky", post.User.Username, post.ID),
		Description: post.Description,
	})
}

func (h *Handlers) renderUser(c *fiber.Ctx) error {
	var (
		username = c.Params("username")
		userID   = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	if username == "" {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Invalid username")
	}

	user, err := h.data.GetUserWithPostsByUsername(userID, username)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching user")
	}

	return utils.RenderPage(c, "user", fiber.Map{"User": user, "ShowFollowButton": userID != user.ID}, &utils.PageDetails{
		Title:       fmt.Sprintf("%s - Stacky", user.Username),
		Description: fmt.Sprintf("%s's stacky profile", user.Username),
	})
}

func (h *Handlers) renderCreatePost(c *fiber.Ctx) error {
	return utils.RenderPage(c, "create", fiber.Map{}, &utils.PageDetails{
		Title:       "Create Post",
		Description: "Create a post on Stacky",
	})
}

func (h *Handlers) renderEditUser(c *fiber.Ctx) error {
	var (
		username = c.Params("username")
		authUser = c.Locals("AuthUser").(*middleware.UserTokenData)
	)

	if username == "" {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Invalid username")
	}

	if username != authUser.Username {
		return utils.RenderError(c, fiber.StatusForbidden, "You do not have permission to edit this user")
	}

	user, err := h.data.GetUserByUsername(authUser.ID, username)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching user")
	}

	return utils.RenderPage(c, "editUser", fiber.Map{"User": user}, &utils.PageDetails{
		Title:       fmt.Sprintf("Edit %s - Stacky", username),
		Description: fmt.Sprintf("Edit %s's stacky profile", username),
	})
}
