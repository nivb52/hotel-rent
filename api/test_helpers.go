package api

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/nivb52/hotel-rent/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DEFAULT_testdb_uri = "mongodb://localhost:27017"
const dbname = "hotel-rent-testing"

type testdb struct {
	db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.Store.User.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func SetupTest(t *testing.T) *testdb {
	err := godotenv.Load("../.env", "../.env.test.local")
	if err != nil {
		log.Fatal("Error loading .env files")
	}
	dburi := os.Getenv("TESTDB_CONNECTION_STRING")
	if len(dburi) < 1 {
		dburi = DEFAULT_testdb_uri
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, dbname)
	store := db.Store{
		User:    db.NewMongoUserStore(client, dbname),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore, dbname),
		Booking: db.NewMongoBookingStore(client, dbname),
	}

	return &testdb{
		Store: store,
	}
}

func (tdb *testdb) afterAll(t *testing.T) {
	defer tdb.teardown(t)
}
