package main

//@Attention: currently seeding is to the main DB - if still command using: db.DBNAME
import (
	"context"
	"fmt"
	"log"

	"github.com/nivb52/hotel-rent/db"
	"github.com/nivb52/hotel-rent/db/fixtures"
	"github.com/nivb52/hotel-rent/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func seedUsers(numberOfUsers int, store *db.Store) int {
	seedUsers := mock.MockUsersMainData(numberOfUsers)
	var errors []error = make([]error, len(*seedUsers))
	for _, user := range *seedUsers {
		_, err := fixtures.AddUser(store, &user, "")
		if err != nil {
			errors = append(errors, err)
		}
	}

	return len(errors)
}

func seedHotels(numberOfHotels int, store *db.Store) {
	seedHotels := mock.MockHotelsMainData(numberOfHotels)

	for _, seedHotel := range seedHotels {
		hotel := mock.MockHotelByInput(&seedHotel)
		numberOfRooms := mock.RandomIntByMaxAndMin(10, 1)
		rooms := mock.MockRoom(numberOfRooms)
		insertedHotel, _ := fixtures.AddHotel(store, &hotel)
		for i := range rooms {
			rooms[i].HotelID = insertedHotel.ID
		}

		updatedCount, _ := fixtures.AddRoom(store, &rooms, insertedHotel.ID.Hex())
		fmt.Printf("Update Hotel with %d rooms \n", updatedCount)
	}
}

func main() {
	//@Todo read from env INITDB_USERNAME & PASSWORD
	mongoURL := db.DBURI
	isDrop := false
	//end env config

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}

	if isDrop {
		if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
			log.Fatal(err)
		}
	}

	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, hotelStore, db.DBNAME)
	userStore := db.NewMongoUserStore(client, db.DBNAME)
	bookingStore := db.NewMongoBookingStore(client, db.DBNAME)
	store := &db.Store{
		User:    userStore,
		Hotel:   hotelStore,
		Room:    roomStore,
		Booking: bookingStore,
	}

	seedHotels(6, store)
	seedUsers(7, store)

}

func init() {
	fmt.Println(" Seeding the db")
}
