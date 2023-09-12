package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nivb52/hotel-rent/db"
)


type UserHandler struct {
	userStore db.UserStore
}


func NewUserHandler(userStore db.UserStore) *UserHandler {
  return &UserHandler {
	userStore: userStore,
  }
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
var( 
	 id = c.Params("id")
)

user, err := h.userStore.GetUserByID(c.Context(), id)
if err != nil {
	return err
}

	return c.JSON(user)
	
}

func HandleGetUsers(c *fiber.Ctx) error {
  return c.SendString("Hello, World HandleGetUsers 👋!")
}