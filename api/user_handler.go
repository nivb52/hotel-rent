package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	e "github.com/nivb52/hotel-rent/api/errors"
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

// Function get a user, returning Json of the a user
func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)

	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusAccepted).JSON(user)
}

func (h *UserHandler) GetUserByEmail(c *fiber.Ctx) error {
	var (
		email = c.Params("email")
	)

	user, err := h.userStore.GetUserByEmail(c.Context(), email)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

// Function get a users, returning Json list of the a users
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, total, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}

	if users == nil {
		return e.ErrResourceNotFound(c)
	}

	if len(users) > 0 {
		return c.JSON(newResourceResp(users, total, 0))
	}

	return e.ErrResourceNotFound(c)
}

// Function create a user and returning Json of the new user
func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var params types.UserParamsForCreate
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(fiber.ErrConflict.Code).JSON(errors)
	}

	userData, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	newUser, err := h.userStore.InsertUser(c.Context(), userData)
	if err != nil {
		fmt.Println("InsertUsere Failed, due: ", err)
		return e.ErrInternalServerError(c)
	}

	return c.Status(fiber.StatusAccepted).JSON(newUser)
}

// Function delete a user and returning Json {msg: "ok", deleted: id}
func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)

	err := h.userStore.DeleteUser(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusAccepted).JSON(map[string]string{"msg": "ok", "deleted": id})
}

// Function update a user and returning Json {msg: "ok", update: id}
func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var params types.UserParamsForUpdate
	if err := c.BodyParser(&params); err != nil {
		return e.ErrBadRequest(c, "Invalid Data (JSON) - failed to parse data")
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	userData, err := types.UpdatedUserFromParams(params)
	if err != nil {
		return err
	}

	_, err = h.userStore.UpdateUserByID(c.Context(), id, userData)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusAccepted).JSON(map[string]string{"msg": "ok", "update": id})
}
