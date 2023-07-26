package handlers

import (
	"strconv"

	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerTodoRoutes(c *fiber.App) {
	r := c.Group("/todos")
	r.Post("/", h.addTodo)
	r.Delete("/:id", h.deleteTodo)
}

func (h *Handlers) addTodo(c *fiber.Ctx) error {
	todo, err := h.data.AddTodo(c.FormValue("todo"))
	if err != nil {
		return utils.SendAlert(c, fiber.StatusBadRequest, "Todo cannot be empty")
	}

	utils.SetAlert(c, fiber.StatusCreated, "Todo added")
	return utils.RenderPartial(c, "todo", todo)
}

func (h *Handlers) deleteTodo(c *fiber.Ctx) error {
	ID, err := strconv.Atoi(c.Params("id", "0"))
	if err != nil {
		return utils.SendAlert(c, fiber.StatusBadRequest, "Invalid ID")
	}

	err = h.data.DeleteTodo(ID)
	if err != nil {
		return utils.SendAlert(c, fiber.StatusBadRequest, "Todo not found")
	}

	utils.SetAlert(c, fiber.StatusNoContent, "Todo deleted")
	return c.SendStatus(fiber.StatusNoContent)
}
