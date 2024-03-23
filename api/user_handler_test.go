package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nivb52/hotel-rent/types"
)

func TestHandleCreateUser(t *testing.T) {
	tdb := SetupTest(t)

	app := fiber.New()
	UserHandler := NewUserHandler(tdb.Store.User)
	app.Post("/", UserHandler.HandleCreateUser)

	params := types.UserParamsForCreate{
		FirstName: "Bob",
		LastName:  "Bob",
		Email:     "bob@google.com",
		Password:  "12345678",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	m := make(map[string]string)
	m["FirstName"] = "expected firstname %s but found %s"
	m["LastName"] = "expected LastName %s but found %s"
	m["Email"] = "expected Email %s but found %s"

	rUser := reflect.ValueOf(user)
	rParams := reflect.ValueOf(params)
	for key, value := range m {
		actual := reflect.Indirect(rUser).FieldByName(key)
		expected := reflect.Indirect(rParams).FieldByName(key)
		if expected.String() != actual.String() {
			t.Errorf(value, expected, actual)
		}
	}

	if user.EncryptedPassword == params.Password {
		t.Errorf("expected Password to not be found but found %s", user.EncryptedPassword)
	}
}
