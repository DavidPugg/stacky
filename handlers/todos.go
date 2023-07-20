package handlers

import (
	"strconv"

	"github.com/davidpugg/stacky/data"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerTodoRoutes(c *fiber.App) {
	c.Get("/", h.renderTodos)
	c.Post("/addTodo", h.addTodo)
	c.Delete("/deleteTodo/:id", h.deleteTodo)
}

func (h *Handlers) renderTodos(c *fiber.Ctx) error {
	todos, _ := h.data.GetTodos()
	return c.Render("index", fiber.Map{
		"Todos": todos,
	})
}

func (h *Handlers) addTodo(c *fiber.Ctx) error {
	todos, _ := h.data.GetTodos()

	var ID int
	if len(todos) == 0 {
		ID = 1
	} else {
		ID = todos[len(todos)-1].ID + 1
	}

	todo := data.Todo{
		ID:   ID,
		Text: c.FormValue("todo"),
	}

	return renderPartial(c, "todo", todo)
}

func (h *Handlers) deleteTodo(c *fiber.Ctx) error {
	ID, _ := strconv.Atoi(c.Params("id", "0"))

	h.data.DeleteTodo(ID)

	return c.Send(nil)
}
