package main

//@Attention: currently seeding is to the main DB - if still command using: db.DBNAME
import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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

		updatedCount, _ := fixtures.AddRooms(store, &rooms, insertedHotel.ID.Hex())
		fmt.Printf("Update Hotel with %d rooms \n", updatedCount)
	}
}

func main() {
	err := godotenv.Load(".env", ".env.test.local")
	if err != nil {
		log.Fatal("Error loading .env files")
	}
	mongoURL := os.Getenv("TESTDB_CONNECTION_STRING")
	mongoDBName := os.Getenv("TESTDB_DATABASE")
	isDropString := os.Getenv("TESTDB_DROP")
	isDrop := false
	if isDropString == "true" || isDropString == "True" || isDropString == "TRUE" {
		isDrop = true
	}

	//end env config

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}

	if mongoDBName == "" {
		mongoDBName = db.DBNAME
	}

	if isDrop {
		if err := client.Database(mongoDBName).Drop(ctx); err != nil {
			log.Fatal(err)
		}
	}

	hotelStore := db.NewMongoHotelStore(client, mongoDBName)
	roomStore := db.NewMongoRoomStore(client, hotelStore, mongoDBName)
	userStore := db.NewMongoUserStore(client, mongoDBName)
	bookingStore := db.NewMongoBookingStore(client, mongoDBName)
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
