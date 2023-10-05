package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nivb52/hotel-rent/db"
	scripts "github.com/nivb52/hotel-rent/scripts/seed"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	hotelStore *db.MongoHotelStore
	roomStore  *db.MongoRoomStore
	userStore  *db.MongoUserStore
	ctx        = context.Background()
)

func main() {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client, db.DBNAME)
	roomStore = db.NewMongoRoomStore(client, hotelStore, db.DBNAME)
	userStore = db.NewMongoUserStore(client, db.DBNAME)

	stores := map[string]interface{}{
		"hotelStore": *db.NewMongoHotelStore(client, db.DBNAME),
		"roomStore":  *db.NewMongoRoomStore(client, hotelStore, db.DBNAME),
		"userStore":  *db.NewMongoUserStore(client, db.DBNAME),
	}

	scripts.SeedHotels(6, ctx, stores)
	scripts.SeedUsers(7, ctx, stores)
}

func init() {
	fmt.Println(" Seeding the db")
}
