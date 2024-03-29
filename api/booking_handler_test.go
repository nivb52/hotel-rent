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
	"github.com/nivb52/hotel-rent/db/fixtures"
	"github.com/nivb52/hotel-rent/types"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetBookingsById(t *testing.T) {
	tdb := SetupTest(t)
	// defer tdb.teardown(t)

	// TODO: refactor to use the mock functions to make it shorter

	//stage
	userData := &types.UserRequiredData{
		Email: "mockEmail@a.com",
		FName: "Alice",
		LName: "Alice",
	}

	insertedUser, err := fixtures.AddUser(&tdb.Store, userData)
	if err != nil {
		log.Fatal(err)
	}

	hotel := types.Hotel{
		Name:     "Grand Hotel",
		Location: "Los Angeles, California",
		Rating:   4,
		CreateAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdateAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	room := types.Room{
		Type:     types.TripleRoomType,
		BedType:  types.TwinBedType,
		Size:     types.RoomSizeKingSize,
		Price:    250,
		CreateAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdateAt: primitive.NewDateTimeFromTime(time.Now()),
	}

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
	m["FromDate"] = TestByField{ErrStr: "expected FromDate to be %s but found %s", ValueType: Date}
	m["ToDate"] = TestByField{ErrStr: "expected ToDate to be %s but found %s", ValueType: Date}
	m["NumPersons"] = TestByField{ErrStr: "expected NumPersons to be %n but found %n", ValueType: Int}
	m["UserID"] = TestByField{ErrStr: "expected UserID to be %s but found %s", ValueType: String}
	m["RoomID"] = TestByField{ErrStr: "expected RoomID to be %s but found %s", ValueType: String}

	assertion := NewDataCompare(m)
	assertion.SetSource(bookingData)
	assertion.SetTestData(resBooking)
	assertion.Compare(t)
}

func TestAdminGetBookings(t *testing.T) {
	tdb := SetupTest(t)
	// defer tdb.teardown(t)

	// TODO: refactor to use the mock functions to make it shorter

	//stage
	userData := &types.UserRequiredData{
		Email: "mockEmail@a.com",
		FName: "Alice",
		LName: "Alice",
	}

	insertedUser, err := fixtures.AddUser(&tdb.Store, userData)
	if err != nil {
		log.Fatal(err)
	}

	hotel := types.Hotel{
		Name:     "Grand Hotel",
		Location: "Los Angeles, California",
		Rating:   4,
		CreateAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdateAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	room := types.Room{
		Type:     types.TripleRoomType,
		BedType:  types.TwinBedType,
		Size:     types.RoomSizeKingSize,
		Price:    250,
		CreateAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdateAt: primitive.NewDateTimeFromTime(time.Now()),
	}

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

	req := httptest.NewRequest("GET", "/", nil)

	// act
	app := fiber.New()
	BookingHandler := NewBookingHandler(&tdb.Store)
	app.Get("/", BookingHandler.AdminGetBookings, nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var resBooking []types.Booking
	json.NewDecoder(resp.Body).Decode(&resBooking)

	// assert
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected to found booking")
	assert.NotNil(t, resBooking)
	assert.GreaterOrEqual(t, len(resBooking), 1)
}
