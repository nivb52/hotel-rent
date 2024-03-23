package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
	"github.com/nivb52/hotel-rent/api/middleware"
	"github.com/nivb52/hotel-rent/db/fixtures"
	"github.com/nivb52/hotel-rent/mock"
	"github.com/nivb52/hotel-rent/types"
)

// # Test Success Auth Request & Token  -  Fn: HandleAuthenticate
func TestSuccessHandleAuthUser(t *testing.T) {
	tdb := SetupTest(t)
	randomInt := mock.RandomIntByMaxAndMin(100, 1)
	randomStr := strconv.Itoa(randomInt)
	mockEmail := "alice" + randomStr + "@google.com"
	pass := "12345678"

	// stage
	userData := &types.UserRequiredData{
		Email: mockEmail,
		FName: "Alice",
		LName: "Alice",
	}

	insertedUser, err := fixtures.AddUser(&tdb.Store, userData, pass)
	if err != nil {
		log.Fatal(err)
	}

	params := AuthParams{
		Email:    mockEmail,
		Password: pass,
	}

	app := fiber.New()
	AuthHandler := NewAuthHandler(tdb.Store.User)
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
		t.Errorf("expected email to be %s - but found %s", insertedUser.Email, claimEmail)
	}

	t.Log(insertedUser.ID.Hex())
	t.Log(insertedUser.ID.String())
	t.Log(insertedUser.ID.MarshalText())
	t.Log(claimId)

	// if claimId != insertedUser.ID.Hex() {
	// 	t.Errorf("expected ID to be %s - but found %s", insertedUser.ID.Hex(), claimId)
	// }

	if parsedToken.Method.Alg() != jwt.SigningMethodHS256.Name {
		t.Errorf("expected algorithm to be %s - but found %s", jwt.SigningMethodES256.Name, parsedToken.Method.Alg())
	}
}

// # Test Failed Auth Request (Wrong Password) - Fn: HandleAuthenticate
func TestFailedHandlAuthUser(t *testing.T) {
	tdb := SetupTest(t)
	// add randomness in the name:
	randomInt := mock.RandomIntByMaxAndMin(100, 1)
	randomStr := strconv.Itoa(randomInt)
	mockEmail := "alice" + randomStr + "@google.com"
	pass := "12345678"

	// stage
	userData := &types.UserRequiredData{
		Email: mockEmail,
		FName: "Alice",
		LName: "Alice",
	}

	_, err := fixtures.AddUser(&tdb.Store, userData, pass)
	if err != nil {
		log.Fatal(err)
	}

	params := AuthParams{
		Email:    mockEmail,
		Password: pass,
	}

	app := fiber.New()
	AuthHandler := NewAuthHandler(tdb.Store.User)
	app.Post("/auth", AuthHandler.HandleAuthenticate)

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	// stage
	wrongParams := AuthParams{
		Email:    mockEmail,
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
		t.Errorf("expected Msg to be <invalid credentials> - but found %s", authBadBody.Msg)
	}
}
