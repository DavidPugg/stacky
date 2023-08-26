package handlers

import (
	"github.com/davidpugg/stacky/internal/data"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Handlers struct {
	mediaEndpoint string
	data          *data.Data
	session       *session.Store
}

func New(data *data.Data, s *session.Store) *Handlers {
	return &Handlers{
		mediaEndpoint: "/media",
		data:          data,
		session:       s,
	}
}

func (h *Handlers) RegisterRoutes(c *fiber.App) {
	h.registerAuthRoutes(c)
	h.registerPostRoutes(c)
	h.registerUserRoutes(c)
	h.registerMediaRoutes(c)
	h.registerViewRoutes(c)
}

func (h *Handlers) getFromSession(c *fiber.Ctx, key string) interface{} {
	session, err := h.session.Get(c)
	if err != nil {
		return nil
	}

	return session.Get(key)
}

type SessionValue struct {
	Key   string
	Value interface{}
}

func (h *Handlers) setSession(c *fiber.Ctx, vals ...SessionValue) error {
	session, err := h.session.Get(c)
	if err != nil {
		return err
	}

	for _, val := range vals {
		session.Set(val.Key, val.Value)
	}

	return session.Save()
}

func (h *Handlers) destoySession(c *fiber.Ctx) error {
	session, err := h.session.Get(c)
	if err != nil {
		return err
	}

	err = session.Destroy()
	if err != nil {
		return err
	}

	return session.Save()
}

func (h *Handlers) getSessionMap(c *fiber.Ctx) map[string]interface{} {
	data := make(map[string]interface{})

	session, err := h.session.Get(c)
	if err != nil {
		return nil
	}

	keys := session.Keys()

	for _, key := range keys {
		data[key] = session.Get(key)
	}

	return data
}
