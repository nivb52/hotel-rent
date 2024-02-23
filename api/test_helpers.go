package api

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/nivb52/hotel-rent/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const testdburi = "mongodb://localhost:27017"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(testdburi))
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
