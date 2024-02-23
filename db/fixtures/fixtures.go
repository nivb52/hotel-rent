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

func AddUser(store *db.Store, input *types.UserRequiredData) error {
	user, err := types.NewUserFromParams(types.UserParamsForCreate{
		Email:     input.Email,
		FirstName: input.FName,
		LastName:  input.LName,
		Password:  "supersecretpassword",
	})

	if err != nil {
		log.Fatal(err)
	}

	_, err = store.User.InsertUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func AddHotel(store *db.Store, hotel *types.Hotel) (*types.Hotel, error) {
	insertedHotel, err := store.Hotel.InsertHotel(ctx, hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New Hotel: ", insertedHotel)
	return hotel, err
}

func AddRoom(store *db.Store, rooms *[]types.Room, hotelID string) (int, error) {
	updatedCount, err := store.Room.InsertRooms(ctx, rooms, hotelID)
	if err != nil {
		log.Fatal(err)
	}

	return updatedCount, err
}
