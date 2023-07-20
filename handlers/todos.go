package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerTodoRoutes(c *fiber.App) {
	c.Get("/", h.renderTodos)
	c.Post("/addTodo", h.addTodo)
	c.Delete("/deleteTodo/:id", h.deleteTodo)
}

func (h *Handlers) renderTodos(c *fiber.Ctx) error {
	todos, _ := h.data.GetTodos()
	return renderPage(c, "index", fiber.Map{"Todos": todos})
}

func (h *Handlers) addTodo(c *fiber.Ctx) error {
	todo, _ := h.data.AddTodo(c.FormValue("todo"))

	return renderPartial(c, "todo", todo)
}

func (h *Handlers) deleteTodo(c *fiber.Ctx) error {
	ID, err := strconv.Atoi(c.Params("id", "0"))
	if err != nil {
		return redirectWithError(c, fiber.StatusNotFound, "Todo not found", h.renderTodos)
	}

	err = h.data.DeleteTodo(ID)
	if err != nil {
		return redirectWithError(c, fiber.StatusNotFound, "Todo not found", h.renderTodos)
	}

	return c.Send(nil)
}
