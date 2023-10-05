package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nivb52/hotel-rent/scripts"
	"github.com/nivb52/hotel-rent/types"
)

func TestHandleAuth(t *testing.T) {
	tdb := setup(t)

	mockUser := scripts.MockUsers(1)[0]
	user, err := types.NewUserFromParams(types.UserParamsForCreate{
		Email:     mockUser.Email,
		FirstName: mockUser.FName,
		LastName:  mockUser.LName,
		Password:  scripts.USER_PASSWORD,
	})

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	insertedUser, err := tdb.UserStore.InsertUser(ctx, user)
	if err != nil {
		fmt.Println("failed to seed user")
		// t.Error("failed to seed data")
		log.Fatal(err)
	}
	fmt.Println("user seedded - success", insertedUser)

	app := fiber.New()

	AuthHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/", AuthHandler.HandleAuthenticate)

	params := AuthParams{
		Password: scripts.USER_PASSWORD,
		Email:    insertedUser.Email,
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)

	var authResp AuthResponse
	json.NewDecoder(resp.Body).Decode(&authResp)
	fmt.Println(authResp)

	// if seededUsers[0].EncryptedPassword == params.Password {
	// 	t.Errorf("expected Passwordto not be found but found %s", seededUsers[0].EncryptedPassword)
	// }
}
