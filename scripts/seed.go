package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
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
	roomStrore := db.NewMongoRoomStore(client, db.DBNAME)

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

		randRoomType, randInt := getRandomRoomType()
		numberOfRooms := randomIntByMaxAndMin(10, 1)

		rooms := []types.Room{
			{
				Type:      randRoomType,
				BasePrice: getRandomPrice(randInt),
				CreateAt:  primitive.NewDateTimeFromTime(time.Now()),
			},
			{
				Type:      randRoomType,
				BasePrice: getRandomPrice(randInt),
				CreateAt:  primitive.NewDateTimeFromTime(time.Now()),
			},
		}
		for i := 0; i < numberOfRooms; i++ {
			rooms = append(rooms, types.Room{
				Type:      randRoomType,
				BasePrice: getRandomPrice(randInt),
				CreateAt:  primitive.NewDateTimeFromTime(time.Now()),
				UpdateAt:  primitive.NewDateTimeFromTime(time.Now()),
			})
		}

		insertedHotel, err := hotelStore.InsertHotel(context.TODO(), &hotel)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("New Hotel: ", insertedHotel)

		for _, room := range rooms {
			room.HotelID = insertedHotel.ID
			insertedRoom, err := roomStrore.InsertRoom(context.TODO(), &room)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("New Room: ", insertedRoom)
		}

	}
}

// Helpers
func getRandomRoomType() (types.RoomType, int) {
	randInt := rand.Intn(int(types.ClosedRoomType))
	switch randInt {
	case 1:
		return types.SingleRoomType, randInt
	case 2:
		return types.DoubleRoomType, randInt
	case 3:
		return types.SeaSideRoomType, randInt
	case 4:
		return types.DeluxRoomType, randInt
	default:
		return types.SingleRoomType, randInt
	}
}

func getRandomPrice(ranInt int) int {
	switch ranInt {
	case 1:
		return randomIntByMaxAndMin(99, 50) //types.SingleRoomType
	case 2:
		return randomIntByMaxAndMin(150, 100) //types.DoubleRoomType
	case 3:
		return randomIntByMaxAndMin(200, 150) //types.SeaSideRoomType
	case 4:
		return randomIntByMaxAndMin(500, 200) //types.DeluxRoomType
	default:
		return randomIntByMaxAndMin(100, 50) //types.SingleRoomType
	}
}

// function which return random integer between to max and min
func randomIntByMaxAndMin(max int, min int) int {
	return rand.Intn(max-min) + min
}
