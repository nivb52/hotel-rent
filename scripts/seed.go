package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nivb52/hotel-rent/db"
	"github.com/nivb52/hotel-rent/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SeedHotel struct {
	Name     string
	Location string
}

func main() {
	fmt.Println(" Seeding the db")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)

	// Create 10 hotel instances
	seedHotels := []SeedHotel{
		{"The Ritz-Carlton", "Los Angeles, California"},
		{"Grand Hyatt", "New York City, New York"},
		{"Marriott Marquis", "Atlanta, Georgia"},
		{"Four Seasons Resort", "Maui, Hawaii"},
		{"Hilton Garden Inn", "Chicago, Illinois"},
		{"Fairmont Empress", "Victoria, British Columbia, Canada"},
		{"The Venetian", "Las Vegas, Nevada"},
		{"Burj Al Arab Jumeirah", "Dubai, United Arab Emirates"},
		{"The Savoy", "London, United Kingdom"},
		{"Hotel del Coronado", "San Diego, California"},
	}

	for _, seedHotel := range seedHotels {
		hotel := types.Hotel{
			Name:     seedHotel.Name,
			Location: seedHotel.Location,
			CreateAt: primitive.NewDateTimeFromTime(time.Now()),
			UpdateAt: primitive.NewDateTimeFromTime(time.Now()),
		}

		insertedHotel, err := hotelStore.InsertHotel(context.TODO(), &hotel)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("New Hotel: ", insertedHotel)
	}
}
