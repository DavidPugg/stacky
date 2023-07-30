package handlers

import (
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerAuthRoutes(c *fiber.App) {
	r := c.Group("/auth")
	r.Post("/login", h.login)
	r.Post("/register", h.register)
	r.Post("/logout", h.logout)
}

func (h *Handlers) login(c *fiber.Ctx) error {
	return c.SendString("login")
}

func (h *Handlers) register(c *fiber.Ctx) error {
	var form struct {
		Avatar string `validate:"required"`
		Username string `validate:"required,min=3,max=32"`
		Email string `validate:"required,email"`
		Password string `validate:"required,min=8,max=32"`
	}

	form.Username = c.FormValue("username")
	form.Email = c.FormValue("email")
	form.Password = c.FormValue("password")
	form.Avatar = "hello.com" //TODO: replace with avatar url
	

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		return utils.SendAlert(c, fiber.StatusBadRequest, "Invalid form data")
	}

	if err := h.data.CreateUser(form.Avatar, form.Username, form.Email, form.Password); err != nil {
		if (err.Error() == "users.username") {
			return utils.SendAlert(c, fiber.StatusBadRequest, "Username already exists")
		}

		if (err.Error() == "users.email") {
			return utils.SendAlert(c, fiber.StatusBadRequest, "Email already exists")
		}
		
		return utils.SendAlert(c, fiber.StatusBadRequest, "Error creating user")
	}
	
	return utils.SendAlert(c, fiber.StatusOK, "User created")
}

func (h *Handlers) logout(c *fiber.Ctx) error {
	return c.SendString("logout")
}