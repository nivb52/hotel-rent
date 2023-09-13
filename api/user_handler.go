package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nivb52/hotel-rent/db"
	"github.com/nivb52/hotel-rent/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)

	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var params types.UserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	userData, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	newUser, err := h.userStore.CreateUser(c.Context(), userData)
	if err != nil {
		return err
	}

	return c.JSON(newUser)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)

	err := h.userStore.DeleteUser(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(map[string]string{"msg": "ok", "deleted": id})
}

// func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx, id string) error {

// }
