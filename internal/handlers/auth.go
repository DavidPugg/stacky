package handlers

import (
	"strings"
	"time"

	"github.com/davidpugg/stacky/internal/middleware"
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerAuthRoutes(c *fiber.App) {
	r := c.Group("/auth")
	r.Post("/validate_email", h.validateEmail)
	r.Post("/validate_username", h.validateUsername)
	r.Post("/validate_password", h.validatePassword)
	r.Post("/login", h.login)
	r.Post("/register", h.register)
	r.Post("/logout", h.logout)
}

func (h *Handlers) validateEmail(c *fiber.Ctx) error {
	var form struct {
		Email string `validate:"required,email"`
	}

	form.Email = c.FormValue("email")

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		if strings.Contains(err.Error(), "email") {
			return c.Status(fiber.StatusBadRequest).SendString("Please enter a valid email")
		}

		if strings.Contains(err.Error(), "required") {
			return c.Status(fiber.StatusBadRequest).SendString("Please enter an email")
		}

		return c.Status(fiber.StatusBadRequest).SendString("Error validating email")
	}

	userID := c.Locals("AuthUser").(*middleware.UserTokenData).ID

	user, err := h.data.GetUserByEmail(userID, form.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("")
	}

	if user != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Email already in use")
	}

	return c.Status(fiber.StatusOK).SendString("")
}

func (h *Handlers) validateUsername(c *fiber.Ctx) error {
	var form struct {
		Username string `validate:"required,min=3,max=32"`
	}

	form.Username = c.FormValue("username")

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		if strings.Contains(err.Error(), "min") || strings.Contains(err.Error(), "max") {
			return c.Status(fiber.StatusBadRequest).SendString("Username must be between 3 and 32 characters")
		}

		if strings.Contains(err.Error(), "required") {
			return c.Status(fiber.StatusBadRequest).SendString("Please enter a username")
		}

		return c.Status(fiber.StatusBadRequest).SendString("Error validating username")
	}

	userID := c.Locals("AuthUser").(*middleware.UserTokenData).ID

	user, err := h.data.GetUserByUsername(userID, form.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("")
	}

	if user != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Username already in use")
	}

	return c.Status(fiber.StatusOK).SendString("")
}

func (h *Handlers) validatePassword(c *fiber.Ctx) error {
	var form struct {
		Password string `validate:"required,min=8,max=32"`
	}

	form.Password = c.FormValue("password")

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		if strings.Contains(err.Error(), "min") || strings.Contains(err.Error(), "max") {
			return c.Status(fiber.StatusBadRequest).SendString("Password must be between 8 and 32 characters")
		}

		if strings.Contains(err.Error(), "required") {
			return c.Status(fiber.StatusBadRequest).SendString("Please enter a password")
		}

		return c.Status(fiber.StatusBadRequest).SendString("Error validating password")
	}

	return c.Status(fiber.StatusOK).SendString("")
}

func (h *Handlers) login(c *fiber.Ctx) error {
	var form struct {
		Username string `validate:"required,min=3,max=32"`
		Password string `validate:"required,min=8,max=32"`
	}

	form.Username = c.FormValue("username")
	form.Password = c.FormValue("password")

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		return utils.SendAlert(c, fiber.StatusBadRequest, "Invalid username or password")
	}

	userID := c.Locals("AuthUser").(*middleware.UserTokenData).ID

	user, err := h.data.GetUserByUsername(userID, form.Username)
	if err != nil {
		return utils.SendAlert(c, fiber.StatusInternalServerError, "Invalid username or password")
	}

	if user == nil {
		return utils.SendAlert(c, fiber.StatusBadRequest, "Invalid username or password")
	}

	if !utils.CheckPasswordHash(form.Password, user.Password) {
		return utils.SendAlert(c, fiber.StatusBadRequest, "Invalid username or password")
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return utils.SendAlert(c, fiber.StatusInternalServerError, "Error generating token")
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	data := middleware.NewUserTokenData(
		user.ID,
		user.Username,
		user.Email,
	)

	utils.SetRedirect(c, "/")
	utils.SetAlert(c, fiber.StatusOK, "Successfully logged in")

	c.Set("HX-Retarget", "#navbar")
	c.Set("HX-Reswap", "outerHTML")
	return utils.RenderPartial(c, "navbar", data)
}

func (h *Handlers) register(c *fiber.Ctx) error {
	var form struct {
		Avatar   string
		Username string `validate:"required,min=3,max=32"`
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=8,max=32"`
	}

	form.Username = c.FormValue("username")
	form.Email = c.FormValue("email")
	form.Password = c.FormValue("password")
	form.Avatar = ""

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		return utils.SendAlert(c, fiber.StatusBadRequest, "Invalid form data")
	}

	hashedPassword, err := utils.HashPassword(form.Password)
	if err != nil {
		return utils.SendAlert(c, fiber.StatusInternalServerError, "Error hashing password")
	}

	if _, err := h.data.CreateUser(form.Avatar, form.Username, form.Email, hashedPassword); err != nil {
		if err.Error() == "users.username" {
			return utils.SendAlert(c, fiber.StatusBadRequest, "Username already exists")
		}

		if err.Error() == "users.email" {
			return utils.SendAlert(c, fiber.StatusBadRequest, "Email already exists")
		}

		return utils.SendAlert(c, fiber.StatusBadRequest, "Error creating user")
	}

	utils.SetRedirect(c, "/login")

	return utils.SendAlert(c, fiber.StatusOK, "User created")
}

func (h *Handlers) logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	c.Locals("AuthUser", nil)

	utils.SetRedirect(c, "/")
	utils.SetAlert(c, fiber.StatusOK, "Successfully logged out")

	c.Set("HX-Retarget", "#navbar")
	c.Set("HX-Reswap", "outerHTML")
	return utils.RenderPartial(c, "navbar", c.Locals("AuthUser"))
}
