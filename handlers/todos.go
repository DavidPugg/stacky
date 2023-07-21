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
	todos, err := h.data.GetTodos()
	if err != nil {
		return renderError(c, fiber.StatusInternalServerError, "Error getting todos")
	}

	return renderPage(c, "index", fiber.Map{"Todos": todos})
}

func (h *Handlers) addTodo(c *fiber.Ctx) error {
	todo, err := h.data.AddTodo(c.FormValue("todo"))
	if err != nil {
		return sendAlert(c, fiber.StatusBadRequest, "Todo cannot be empty")
	}

	setAlert(c, fiber.StatusCreated, "Todo added")
	return renderPartial(c, "todo", todo)
}

func (h *Handlers) deleteTodo(c *fiber.Ctx) error {
	ID, err := strconv.Atoi(c.Params("id", "0"))
	if err != nil {
		return sendAlert(c, fiber.StatusBadRequest, "Invalid ID")
	}

	err = h.data.DeleteTodo(ID)
	if err != nil {
		return sendAlert(c, fiber.StatusBadRequest, "Todo not found")
	}

	setAlert(c, fiber.StatusOK, "Todo deleted")
	return c.Send(nil)
}
