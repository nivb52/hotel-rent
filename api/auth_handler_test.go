package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
	"github.com/nivb52/hotel-rent/api/middleware"
	"github.com/nivb52/hotel-rent/db"
	"github.com/nivb52/hotel-rent/types"
)

// END COMMON
func insertUser(userStore db.UserStore, userParams types.UserParamsForCreate) (*types.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userData, err := types.NewUserFromParams(userParams)

	if err != nil {
		return nil, err
	}

	user, err := userStore.InsertUser(ctx, userData)
	if err != nil {
		return nil, err
	}
	fmt.Println("seed user created in db")
	return user, nil
}

func TestHandlAuthUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	// # Test (Success) Auth Request & Token

	// stage
	userData := types.UserParamsForCreate{
		FirstName: "Bob",
		LastName:  "Alice",
		Email:     "alice@google.com",
		Password:  "12345678",
	}

	insertedUser, err := insertUser(tdb.UserStore, userData)
	if err != nil {
		log.Fatal(err)
	}

	params := AuthParams{
		Email:    "alice@google.com",
		Password: "12345678",
	}

	app := fiber.New()
	AuthHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", AuthHandler.HandleAuthenticate)

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	// act
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var authBody AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authBody)
	if err != nil {
		log.Println("failed to decode response")
		t.Error(err)
	}

	// expect (Auth)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected HTTP Status OK - but found %s", resp.Status)
	}

	if authBody.Token == "" {
		t.Errorf("expected Token exists - but found %s", authBody.Token)
	}

	if authBody.User.Email != params.Email {
		t.Errorf("expected email to be %s - but found %s", params.Email, authBody.User.Email)
	}

	// expect (Token)
	secretKey := middleware.GetSecret()
	parsedToken, err := jwt.Parse(authBody.Token, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	claims := parsedToken.Claims.(jwt.MapClaims)
	claimEmail := claims["email"].(string)
	claimId := claims["id"].(string)

	if claimEmail != insertedUser.Email {
		t.Errorf("expected algorithm to be %s - but found %s", insertedUser.Email, claimEmail)
	}

	if claimId != insertedUser.ID.Hex() {
		t.Errorf("expected algorithm to be %s - but found %s", insertedUser.ID.Hex(), claimId)
	}

	if parsedToken.Method.Alg() != jwt.SigningMethodHS256.Name {
		t.Errorf("expected algorithm to be %s - but found %s", jwt.SigningMethodES256.Name, parsedToken.Method.Alg())
	}

	// # Test Wrong Password (Failed Auth Request)
	// stage
	wrongParams := AuthParams{
		Email:    "alice@google.com",
		Password: "WORNG_PASSWORD",
	}

	by, _ := json.Marshal(wrongParams)
	badReq := httptest.NewRequest("POST", "/auth", bytes.NewReader(by))
	badReq.Header.Add("Content-Type", "application/json")

	// act
	bedResp, err := app.Test(badReq)
	if err != nil {
		t.Error(err)
	}
	var authBadBody ForbiddenResponse
	err = json.NewDecoder(bedResp.Body).Decode(&authBadBody)
	if err != nil {
		log.Println("failed to decode response")
		t.Error(err)
	}

	// expect
	if bedResp.StatusCode != http.StatusForbidden {
		t.Errorf("expected HTTP Status Forbidden - but found %s", bedResp.Status)
	}

	if authBadBody.Msg != "invalid credentials" {
		t.Errorf("expected Token to no exists - but found %s", authBadBody.Msg)
	}
}
