package handlers

import (
	"fmt"
	"strconv"

	"github.com/davidpugg/stacky/internal/data"
	"github.com/davidpugg/stacky/internal/middleware"
	"github.com/davidpugg/stacky/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) registerViewRoutes(c *fiber.App) {

	c.Use(func(c *fiber.Ctx) error {
		authUserID := c.Locals("AuthUser").(*middleware.UserTokenData).ID

		if c.Get("HX-Request") != "true" {
			followedUsers, err := h.data.GetAllFollowedUsers(authUserID)
			if err != nil {
				return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching followed users")
			}

			c.Locals("FollowedUsers", followedUsers)
		}

		return c.Next()
	})

	c.Get("/", middleware.MainAuthGuard, h.renderMain)
	c.Get("/discover", h.renderDiscover)
	c.Get("/register", h.renderRegister)
	c.Get("/login", middleware.LoginAuthGuard, h.renderLogin)
	c.Get("/post/:id", h.renderPost)
	c.Get("/u/:username", h.renderUser)
	c.Get("/create", middleware.MainAuthGuard, h.renderCreatePost)
	c.Get("/u/:username/edit", middleware.Authenticate, h.renderEditUser)
}

func (h *Handlers) renderMain(c *fiber.Ctx) error {
	var (
		authUserID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	posts, err := h.data.GetFollowedPosts(authUserID, 1, postsPageLimit)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching posts")
	}

	rp := utils.CreateRevealeObjects(posts, 1, postsPageLimit)

	return utils.RenderPage(c, "index", fiber.Map{"Posts": rp}, &utils.PageDetails{
		Title:       "Stacky",
		Description: "Stacky is a simple social media platform",
	})
}

func (h *Handlers) renderDiscover(c *fiber.Ctx) error {
	var (
		authUserID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
	)

	posts, err := h.data.GetAllPosts(authUserID, 1, postsPageLimit)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching posts")
	}

	rp := utils.CreateRevealeObjects(posts, 1, postsPageLimit)

	return utils.RenderPage(c, "index", fiber.Map{"Posts": rp}, &utils.PageDetails{
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
	var (
		pID         = c.Params("id")
		authUserID  = c.Locals("AuthUser").(*middleware.UserTokenData).ID
		commentChan = make(chan []*data.Comment)
		errorChan   = make(chan error)
	)

	if pID == "" {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Invalid post ID")
	}

	postID, err := strconv.Atoi(pID)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Invalid post ID")
	}

	go func() {
		comments, err := h.data.GetPostComments(authUserID, postID, 1, commentsPageLimit)
		if err != nil {
			errorChan <- err
			return
		}

		errorChan <- nil
		commentChan <- comments
	}()

	post, err := h.data.GetPostByID(authUserID, postID)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching posts")
	}

	err = <-errorChan
	if err != nil {
		return err
	}

	comments := <-commentChan
	rc := utils.CreateRevealeObjects(comments, 1, smallPostsPageLimit)

	return utils.RenderPage(c, "post", fiber.Map{"Post": post, "Comments": rc, "ShowFollowButton": authUserID != post.User.ID}, &utils.PageDetails{
		Title:       fmt.Sprintf("%s - %d - Stacky", post.User.Username, post.ID),
		Description: post.Description,
	})
}

func (h *Handlers) renderUser(c *fiber.Ctx) error {
	var (
		username   = c.Params("username")
		authUserID = c.Locals("AuthUser").(*middleware.UserTokenData).ID
		postsChan  = make(chan []*data.Post)
		errorChan  = make(chan error)
	)

	if username == "" {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Invalid username")
	}

	go func() {
		posts, err := h.data.GetPostsOfUserByUsername(authUserID, username, 1, smallPostsPageLimit)
		if err != nil {
			errorChan <- err
			return
		}

		errorChan <- nil
		postsChan <- posts
	}()

	user, err := h.data.GetUserByUsername(authUserID, username)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching user")
	}

	err = <-errorChan
	if err != nil {
		return err
	}

	posts := <-postsChan
	rp := utils.CreateRevealeObjects(posts, 1, smallPostsPageLimit)

	return utils.RenderPage(c, "user", fiber.Map{"User": user, "Posts": rp, "ShowFollowButton": authUserID != user.ID}, &utils.PageDetails{
		Title:       fmt.Sprintf("%s - Stacky", user.Username),
		Description: fmt.Sprintf("%s's stacky profile", user.Username),
	})
}

func (h *Handlers) renderCreatePost(c *fiber.Ctx) error {
	return utils.RenderPage(c, "create", fiber.Map{}, &utils.PageDetails{
		Title:       "Create Post",
		Description: "Create a post on Stacky",
	})
}

func (h *Handlers) renderEditUser(c *fiber.Ctx) error {
	var (
		username = c.Params("username")
		authUser = c.Locals("AuthUser").(*middleware.UserTokenData)
	)

	if username == "" {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Invalid username")
	}

	if username != authUser.Username {
		return utils.RenderError(c, fiber.StatusForbidden, "You do not have permission to edit this user")
	}

	user, err := h.data.GetUserByUsername(authUser.ID, username)
	if err != nil {
		return utils.RenderError(c, fiber.StatusInternalServerError, "Error fetching user")
	}

	return utils.RenderPage(c, "editUser", fiber.Map{"User": user}, &utils.PageDetails{
		Title:       fmt.Sprintf("Edit %s - Stacky", username),
		Description: fmt.Sprintf("Edit %s's stacky profile", username),
	})
}
