package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nivb52/hotel-rent/api/middleware"
	"github.com/nivb52/hotel-rent/db/fixtures"
	"github.com/nivb52/hotel-rent/mock"
	"github.com/nivb52/hotel-rent/types"
	"github.com/stretchr/testify/assert"
)

// TODO: refactor to use the mock functions to make it shorter

func TestGetBookingsById(t *testing.T) {
	tdb := SetupTest(t)
	// defer tdb.teardown(t)

	//stage
	user := mock.User()
	insertedUser, err := fixtures.AddUser(&tdb.Store, &user)
	if err != nil {
		log.Fatal(err)
	}
	hotel := mock.Hotel()
	room := mock.MockRoom(1)[0]

	insertedHotel, _ := fixtures.AddHotel(&tdb.Store, &hotel)
	room.HotelID = insertedHotel.ID
	rooms := []types.Room{}
	rooms = append(rooms, room)
	insertedRoom, err := fixtures.AddRoom(&tdb.Store, &room)
	if err != nil {
		log.Fatal(err)
	}

	from := time.Now()
	till := from.AddDate(0, 0, 5)
	params := types.BookingParamsForCreate{
		UserID:     insertedUser.ID.Hex(),
		RoomID:     insertedRoom.ID.Hex(),
		FromDate:   from,
		TillDate:   till,
		NumPersons: 4,
	}

	bookingData, err := fixtures.AddBooking(&tdb.Store, &params)
	if err != nil {
		log.Fatal(err)
	}

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

	//stage - insert into database
	//stage
	user := mock.User()
	insertedUser, err := fixtures.AddUser(&tdb.Store, &user)
	if err != nil {
		log.Fatal(err)
	}
	hotel := mock.Hotel()
	room := mock.MockRoom(1)[0]

	insertedHotel, _ := fixtures.AddHotel(&tdb.Store, &hotel)
	room.HotelID = insertedHotel.ID
	rooms := []types.Room{}
	rooms = append(rooms, room)
	insertedRoom, err := fixtures.AddRoom(&tdb.Store, &room)
	if err != nil {
		log.Fatal(err)
	}

	from := time.Now()
	till := from.AddDate(0, 0, 5)
	params := types.BookingParamsForCreate{
		UserID:     insertedUser.ID.Hex(),
		RoomID:     insertedRoom.ID.Hex(),
		FromDate:   from,
		TillDate:   till,
		NumPersons: 4,
	}

	_, err = fixtures.AddBooking(&tdb.Store, &params)
	if err != nil {
		log.Fatal(err)
	}

	err = fixtures.SetAdminUser(&tdb.Store, insertedUser)
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
	user := mock.User()
	insertedUser, err := fixtures.AddUser(&tdb.Store, &user)
	if err != nil {
		log.Fatal(err)
	}
	hotel := mock.Hotel()
	room := mock.MockRoom(1)[0]

	insertedHotel, _ := fixtures.AddHotel(&tdb.Store, &hotel)
	room.HotelID = insertedHotel.ID
	rooms := []types.Room{}
	rooms = append(rooms, room)
	insertedRoom, err := fixtures.AddRoom(&tdb.Store, &room)
	if err != nil {
		log.Fatal(err)
	}

	from := time.Now()
	till := from.AddDate(0, 0, 5)
	params := types.BookingParamsForCreate{
		UserID:     insertedUser.ID.Hex(),
		RoomID:     insertedRoom.ID.Hex(),
		FromDate:   from,
		TillDate:   till,
		NumPersons: 4,
	}

	_, err = fixtures.AddBooking(&tdb.Store, &params)
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
	assert.NotEqual(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, resBooking)
	assert.Equal(t, len(resBooking), 0)
}
