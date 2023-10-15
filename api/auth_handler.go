package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nivb52/hotel-rent/api/middleware"
	"github.com/nivb52/hotel-rent/db"
	"github.com/nivb52/hotel-rent/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user,omitempty"`
	Token string      `json:"token,omitempty"`
}

type ForbiddenResponse struct {
	Msg    string `json:"msg"`
	Reason string `json:"reason"`
}

// return ForbiddenRequest
func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(ForbiddenResponse{
		Msg:    "invalid credentials", // fiber.ErrForbidden.Message
		Reason: "wrong password or email",
	})
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	err := c.BodyParser(&params)
	if err != nil {
		fmt.Println("Failed BodyParser")
		return err
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}

	isValid := types.IsPasswordValid(user.EncryptedPassword, params.Password)
	if !isValid {
		return invalidCredentials(c)
	}

	tokenString, err := createToken(user.ID.Hex(), user.Email, user.IsAdmin)
	if err != nil {
		return err
	}

	fmt.Println("Authenticated user")
	resp := AuthResponse{
		User:  user,
		Token: tokenString,
	}
	return c.JSON(resp)
}

// create token From User data of id and email.
func createToken(userID, userEmail string, isAdmin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":      userID,
			"email":   userEmail,
			"isAdmin": isAdmin,
			"exp":     time.Now().Add(time.Hour * 12).UTC().Unix(),
		})

	secret := middleware.GetSecret()
	tokenString, err := token.SignedString(secret)
	if err != nil {
		fmt.Println("Failed to sign token with secret, userID: ", userID, " useEmail: ", userEmail)
		return "", err
	}

	return tokenString, nil
}
