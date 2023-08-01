package handlers

import (
	"strconv"

	"github.com/davidpugg/stacky/internal/data"
	"github.com/davidpugg/stacky/internal/middleware"
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerViewRoutes(c *fiber.App) {
	c.Get("/", middleware.UpdatePageDetails, h.renderMain)
	c.Get("/login", middleware.UpdatePageDetails, h.renderLogin)
	c.Get("/register", middleware.UpdatePageDetails, h.renderRegister)
	c.Get("/post/:id", middleware.UpdatePageDetails, h.renderPost)
}

func (h *Handlers) renderMain(c *fiber.Ctx) error {
	posts, err := h.data.GetPosts()
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching posts")
	}

	return utils.RenderPage(c, "index", fiber.Map{"Posts": posts})
}

func (h *Handlers) renderLogin(c *fiber.Ctx) error {
	return utils.RenderPage(c, "login", fiber.Map{})
}

func (h *Handlers) renderRegister(c *fiber.Ctx) error {
	return utils.RenderPage(c, "register", fiber.Map{})
}

func (h *Handlers) renderPost(c *fiber.Ctx) error {
	posts, err := h.data.GetPosts()
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching posts")
	}

	postID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error converting post id to int")
	}

	var post *data.Post

	for _, p := range posts {
		if p.ID == postID {
			post = p
			break
		}
	}

	return utils.RenderPage(c, "post", fiber.Map{"Post": post})
}
