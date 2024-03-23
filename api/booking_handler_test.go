package api

import (
	"log"
	"testing"
	"time"

	"github.com/nivb52/hotel-rent/db/fixtures"
	"github.com/nivb52/hotel-rent/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetBookings(t *testing.T) {
	tdb := SetupTest(t)
	// defer db.teardown(t)

	// TODO: refactor to use the mock functions

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

	// act
	from := time.Now()
	till := from.AddDate(0, 0, 5)
	params := types.BookingParamsForCreate{
		UserID:     insertedUser.ID.Hex(),
		RoomID:     insertedRoom.ID.Hex(),
		FromDate:   from,
		TillDate:   till,
		NumPersons: 4,
	}

	booking, err := fixtures.AddBooking(&tdb.Store, &params)
	if err != nil {
		log.Fatal(err)
	}

	// assert
	if booking.FromDate != params.FromDate {
		formattedParams := params.FromDate.Format("01/02/2006")
		formattedBooking := booking.FromDate.Format("01/02/2006")
		t.Errorf("expected from data to be %s but found %s", formattedParams, formattedBooking)
	}

	if booking.NumPersons != params.NumPersons {
		t.Errorf("expected Persons to be %d found but found %d", params.NumPersons, booking.NumPersons)
	}

}
