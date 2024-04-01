package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nivb52/hotel-rent/api/middleware"
	"github.com/nivb52/hotel-rent/db/fixtures"
	"github.com/nivb52/hotel-rent/mock"
	"github.com/nivb52/hotel-rent/types"
	"github.com/stretchr/testify/assert"
)

func insertBooking(tdb *testdb) (*types.Booking, *types.User) {
	user := mock.User()
	insertedUser, err := fixtures.AddUser(&tdb.Store, &user)
	if err != nil {
		log.Fatal(err)
	}
	hotel := mock.Hotel()
	room := mock.MockRoom(1)[0]

	insertedHotel, _ := fixtures.AddHotel(&tdb.Store, &hotel)
	room.HotelID = insertedHotel.ID
	insertedRoom, err := fixtures.AddRoom(&tdb.Store, &room)
	if err != nil {
		log.Fatal(err)
	}

	params := mock.Booking(insertedUser.ID.Hex(), insertedRoom.ID.Hex())
	insertedBooking, err := fixtures.AddBooking(&tdb.Store, &params)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking, insertedUser
}

func TestGetBookingsById(t *testing.T) {
	tdb := SetupTest(t)
	// defer tdb.teardown(t)

	//stage
	bookingData, _ := insertBooking(tdb)
	id := bookingData.ID.Hex()
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", id), nil)

	// act
	app := fiber.New()
	BookingHandler := NewBookingHandler(&tdb.Store)
	// 		Make a request to the endpoint with a specific id
	app.Get("/:id", BookingHandler.GetBookingsById, nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var resBooking types.Booking
	json.NewDecoder(resp.Body).Decode(&resBooking)

	// assert
	assert.Equal(t, http.StatusOK, resp.StatusCode, "No Booking Found - Expected to found booking")

	// (prepare for assersion )
	m := make(map[string]TestByField)
	m["FromDate"] = TestByField{ErrStr: "", ValueType: Date}
	m["ToDate"] = TestByField{ErrStr: "", ValueType: Date}
	m["NumPersons"] = TestByField{ErrStr: "", ValueType: Int}
	m["UserID"] = TestByField{ErrStr: "", ValueType: String}
	m["RoomID"] = TestByField{ErrStr: "", ValueType: String}

	assertion := NewDataCompare(m)
	assertion.SetSource(bookingData)
	assertion.SetTestData(resBooking)
	assertion.Compare(t)
}

func TestAdminGetBookings(t *testing.T) {
	tdb := SetupTest(t)
	// defer tdb.teardown(t)

	//stage
	_, insertedUser := insertBooking(tdb)
	err := fixtures.SetAdminUser(&tdb.Store, insertedUser)
	if err != nil {
		log.Fatal(err)
	}

	BookingHandler := NewBookingHandler(&tdb.Store)
	app := fiber.New()
	adminRoute := app.Group("/")
	adminRoute.Get("/", middleware.JWTAuthentication, middleware.IsAdminAuth, BookingHandler.AdminGetBookings, nil)

	token, err := createToken(insertedUser.ID.String(), insertedUser.Email, insertedUser.IsAdmin)
	if err != nil {
		log.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", token)

	// act
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var resBooking []types.Booking
	json.NewDecoder(resp.Body).Decode(&resBooking)

	// assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, resBooking)
	assert.GreaterOrEqual(t, len(resBooking), 1)
}

// test non admin cannot aceess the bookings
func TestAdminGetBookingsWithNonAdmin(t *testing.T) {
	tdb := SetupTest(t)
	// defer tdb.teardown(t)

	//stage
	_, insertedUser := insertBooking(tdb)
	BookingHandler := NewBookingHandler(&tdb.Store)
	app := fiber.New()
	adminRoute := app.Group("/")
	adminRoute.Get("/", middleware.JWTAuthentication, middleware.IsAdminAuth, BookingHandler.AdminGetBookings, nil)

	token, err := createToken(insertedUser.ID.String(), insertedUser.Email, insertedUser.IsAdmin)
	if err != nil {
		log.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", token)

	// act
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var resBooking []types.Booking
	json.NewDecoder(resp.Body).Decode(&resBooking)

	// assert
	assert.NotEqual(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, resBooking)
	assert.Equal(t, len(resBooking), 0)
}
