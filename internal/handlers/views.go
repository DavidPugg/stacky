package handlers

import (
	"fmt"

	"github.com/davidpugg/stacky/internal/middleware"
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerViewRoutes(c *fiber.App) {
	c.Get("/", h.renderMain)
	c.Get("/login", h.renderLogin)
	c.Get("/register", h.renderRegister)
	c.Get("/post/:id", h.renderPost)
}

func (h *Handlers) renderMain(c *fiber.Ctx) error {
	userID := c.Locals("User").(*middleware.UserTokenData).ID

	posts, err := h.data.GetPosts(userID)
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
	postID := c.Params("id")
	if postID == "" {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Invalid post ID")
	}

	userID := c.Locals("User").(*middleware.UserTokenData).ID

	post, err := h.data.GetPostByID(userID, postID)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching posts")
	}
	return utils.RenderPage(c, "post", fiber.Map{"Post": post}, &utils.PageDetails{
		Title:       fmt.Sprintf("%s - %d - Stacky", post.User.Username, post.ID),
		Description: post.Description,
	})
}
