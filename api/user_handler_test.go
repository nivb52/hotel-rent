package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nivb52/hotel-rent/db"
	"github.com/nivb52/hotel-rent/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const testdburi = "mongodb://localhost:27017"
const dbname = "hotel-rent-testing"

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client, dbname),
	}
}

func (tdb *testdb) afterAll(t *testing.T) {
	defer tdb.teardown(t)
}

// END COMMON

func TestHandleCreateUser(t *testing.T) {
	tdb := setup(t)

	app := fiber.New()
	UserHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", UserHandler.HandleCreateUser)

	params := types.UserParamsForCreate{
		FirstName: "Bob",
		LastName:  "Alice",
		Email:     "alice@google.com",
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
		t.Errorf("expected Passwordto not be found but found %s", user.EncryptedPassword)
	}
}
