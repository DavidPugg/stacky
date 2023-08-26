package handlers

import (
	"strings"

	"github.com/davidpugg/stacky/internal/middleware"
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerAuthRoutes(c *fiber.App) {
	r := c.Group("/auth")

	r.Post("/login", h.login)
	r.Post("/logout", h.logout)
	r.Post("/register", h.register)

	r.Post("/validate_email", h.validateEmail)
	r.Post("/validate_username", h.validateUsername)
	r.Post("/validate_password", h.validatePassword)
}

type Form struct {
	Username string `validate:"required,min=3,max=14"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=32"`
}

func (h *Handlers) validateEmail(c *fiber.Ctx) error {
	var (
		form   = Form{Email: c.FormValue("email")}
		userID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	validate := validator.New()
	if err := validate.StructPartial(form, "Email"); err != nil {
		if strings.Contains(err.Error(), "required") {
			return c.Status(fiber.StatusBadRequest).SendString("Please enter an email")
		}

		return c.Status(fiber.StatusBadRequest).SendString("Please enter a valid email")
	}

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
	var (
		form   = Form{Username: c.FormValue("username")}
		userID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	validate := validator.New()
	if err := validate.StructPartial(form, "Username"); err != nil {
		if strings.Contains(err.Error(), "min") || strings.Contains(err.Error(), "max") {
			return c.Status(fiber.StatusBadRequest).SendString("Username must be between 3 and 32 characters")
		}

		if strings.Contains(err.Error(), "required") {
			return c.Status(fiber.StatusBadRequest).SendString("Please enter a username")
		}

		return c.Status(fiber.StatusBadRequest).SendString("Error validating username")
	}

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
	var (
		form = Form{Password: c.FormValue("password")}
	)

	validate := validator.New()
	if err := validate.StructPartial(form, "Password"); err != nil {
		if strings.Contains(err.Error(), "min") || strings.Contains(err.Error(), "max") {
			return c.Status(fiber.StatusBadRequest).SendString("Password must be between 8 and 32 characters")
		}

		return c.Status(fiber.StatusBadRequest).SendString("Please enter a valid password")
	}

	return c.Status(fiber.StatusOK).SendString("")
}

func (h *Handlers) login(c *fiber.Ctx) error {
	var (
		form = Form{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}
		userID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	validate := validator.New()
	if err := validate.StructPartial(form, "Username", "Password"); err != nil {
		return utils.SendAlert(c, fiber.StatusBadRequest, "Invalid username or password")
	}

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

	session, err := h.session.Get(c)
	if err != nil {
		return utils.SendAlert(c, fiber.StatusInternalServerError, "Error getting session")
	}

	session.Set("id", user.ID)
	session.Set("name", user.Username)
	session.Set("email", user.Email)
	session.Set("avatar", user.Avatar)

	if err := session.Save(); err != nil {
		panic(err)
	}

	data := middleware.NewUserTokenData(
		user.ID,
		user.Avatar,
		user.Username,
		user.Email,
	)

	utils.SetRedirect(c, "/")
	utils.SetAlert(c, fiber.StatusOK, "Successfully logged in")

	return utils.RenderPartial(c, "navbar", data)
}

func (h *Handlers) register(c *fiber.Ctx) error {
	var (
		form = Form{
			Username: c.FormValue("username"),
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
		}
	)

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		return utils.SendAlert(c, fiber.StatusBadRequest, "Invalid form data")
	}

	hashedPassword, err := utils.HashPassword(form.Password)
	if err != nil {
		return utils.SendAlert(c, fiber.StatusInternalServerError, "Error hashing password")
	}

	if _, err := h.data.CreateUser("", form.Username, form.Email, hashedPassword); err != nil {
		if err.Error() == "username_unique" {
			return utils.SendAlert(c, fiber.StatusBadRequest, "Username already exists")
		}

		if err.Error() == "email_unique" {
			return utils.SendAlert(c, fiber.StatusBadRequest, "Email already exists")
		}

		return utils.SendAlert(c, fiber.StatusBadRequest, "Error creating user")
	}

	utils.SetRedirect(c, "/login")

	return utils.SendAlert(c, fiber.StatusOK, "User created")
}

func (h *Handlers) logout(c *fiber.Ctx) error {
	session, err := h.session.Get(c)
	if err != nil {
		return utils.SendAlert(c, fiber.StatusInternalServerError, "Error getting session")
	}

	if err := session.Destroy(); err != nil {
		return utils.SendAlert(c, fiber.StatusInternalServerError, "Error destroying session")
	}

	if err := session.Save(); err != nil {
		panic(err)
	}

	utils.SetRedirect(c, "/")
	utils.SetAlert(c, fiber.StatusOK, "Successfully logged out")

	return utils.RenderPartial(c, "navbar", middleware.UserTokenData{})
}
