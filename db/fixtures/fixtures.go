package fixtures

import (
	"context"
	"fmt"
	"log"

	"github.com/nivb52/hotel-rent/db"
	"github.com/nivb52/hotel-rent/types"
)

var (
	ctx = context.Background()
)

func AddUser(store *db.Store, input *types.UserRequiredData, overridePass ...string) (*types.User, error) {
	var pass string
	if len(overridePass) < 1 {
		pass = "supersecretpassword"
	} else {
		pass = overridePass[0]
	}

	user, err := types.NewUserFromParams(types.UserParamsForCreate{
		Email:     input.Email,
		FirstName: input.FName,
		LastName:  input.LName,
		Password:  pass,
	})

	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := store.User.InsertUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return insertedUser, nil
}

func AddHotel(store *db.Store, hotel *types.Hotel) (*types.Hotel, error) {
	insertedHotel, err := store.Hotel.InsertHotel(ctx, hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New Hotel: ", insertedHotel)
	return hotel, err
}

func AddRooms(store *db.Store, rooms *[]types.Room, hotelID string) (int, error) {
	updatedCount, err := store.Room.InsertRooms(ctx, rooms, hotelID)
	if err != nil {
		log.Fatal(err)
	}

	return updatedCount, err
}

func AddRoom(store *db.Store, room *types.Room) (*types.Room, error) {
	insertedRoom, err := store.Room.InsertRoom(ctx, room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom, err
}

func AddBooking(store *db.Store, booking *types.BookingParamsForCreate) (*types.Booking, error) {
	insertedBooking, err := store.Booking.InsertBooking(ctx, booking)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking, err
}

func SetAdminUser(store *db.Store, user *types.User) error {
	user.IsAdmin = true
	_, err := store.User.UpdateUserByID(ctx, user.ID.Hex(), user)
	if err != nil {
		return err
	}

	return nil
}
