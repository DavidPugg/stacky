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
	c.Get("/login", middleware.LoginAuthGuard, h.renderLogin)
	c.Get("/register", h.renderRegister)
	c.Get("/post/:id", h.renderPost)
	c.Get("/discover", h.renderDiscover)
	c.Get("/u/:username", h.renderUser)
	c.Get("/create", middleware.MainAuthGuard, h.renderCreatePost)
}

func (h *Handlers) renderMain(c *fiber.Ctx) error {
	userID := c.Locals("AuthUser").(*middleware.UserTokenData).ID

	posts, err := h.data.GetFollowedPosts(userID)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching posts")
	}

	return utils.RenderPage(c, "index", fiber.Map{"Posts": posts}, &utils.PageDetails{
		Title:       "Stacky",
		Description: "Stacky is a simple social media platform",
	})
}

func (h *Handlers) renderDiscover(c *fiber.Ctx) error {
	userID := c.Locals("AuthUser").(*middleware.UserTokenData).ID

	posts, err := h.data.GetAllPosts(userID)
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
	pID := c.Params("id")
	if pID == "" {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Invalid post ID")
	}

	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Invalid post ID")
	}

	userID := c.Locals("AuthUser").(*middleware.UserTokenData).ID

	post, err := h.data.GetPostWithCommentsByID(userID, postID)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching posts")
	}

	return utils.RenderPage(c, "post", fiber.Map{"Post": post}, &utils.PageDetails{
		Title:       fmt.Sprintf("%s - %d - Stacky", post.User.Username, post.ID),
		Description: post.Description,
	})
}

func (h *Handlers) renderUser(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Invalid username")
	}

	userID := c.Locals("AuthUser").(*middleware.UserTokenData).ID

	user, err := h.data.GetUserWithPostsByUsername(userID, username)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching user")
	}

	return utils.RenderPage(c, "user", fiber.Map{"User": user}, &utils.PageDetails{
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
