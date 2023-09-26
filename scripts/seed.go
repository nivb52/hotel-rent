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
	Rating   int8
}

var (
	client     *mongo.Client
	hotelStore *db.MongoHotelStore
	roomStore  *db.MongoRoomStore
	ctx        = context.Background()
)

func seedHotels(numberOfHotels int) {

	// 10 hotel options
	seedHotels10 := []SeedHotel{
		{"The Ritz-Carlton", "Los Angeles, California", 5},
		{"Grand Hyatt", "New York City, New York", 4},
		{"Marriott Marquis", "Atlanta, Georgia", 3},
		{"Four Seasons Resort", "Maui, Hawaii", 5},
		{"Hilton Garden Inn", "Chicago, Illinois", 3},
		{"Fairmont Empress", "Victoria, British Columbia, Canada", 5},
		{"The Venetian", "Las Vegas, Nevada", 4},
		{"Burj Al Arab Jumeirah", "Dubai, United Arab Emirates", 5},
		{"The Savoy", "London, United Kingdom", 3},
		{"Hotel del Coronado", "San Diego, California", 4},
	}

	seedHotels := seedHotels10[0:numberOfHotels]

	for _, seedHotel := range seedHotels {
		hotel := types.Hotel{
			Name:     seedHotel.Name,
			Location: seedHotel.Location,
			Rating:   seedHotel.Rating,
			CreateAt: primitive.NewDateTimeFromTime(time.Now()),
			UpdateAt: primitive.NewDateTimeFromTime(time.Now()),
		}

		randRoomType, randInt := getRandomRoomType()
		randBedType := getRandomBedType(types.RoomType(randInt))
		randSize := getRandomSizeStirng(types.RoomType(randInt))
		numberOfRooms := randomIntByMaxAndMin(10, 1)

		rooms := []types.Room{
			{
				Type:     randRoomType,
				BedType:  randBedType,
				Size:     randSize,
				Price:    getRandomPrice(randInt),
				CreateAt: primitive.NewDateTimeFromTime(time.Now()),
				UpdateAt: primitive.NewDateTimeFromTime(time.Now()),
			},
			{
				Type:     randRoomType,
				BedType:  randBedType,
				Size:     randSize,
				Price:    getRandomPrice(randInt),
				CreateAt: primitive.NewDateTimeFromTime(time.Now()),
				UpdateAt: primitive.NewDateTimeFromTime(time.Now()),
			},
		}

		for i := 0; i < numberOfRooms; i++ {
			rooms = append(rooms, types.Room{
				Type:     randRoomType,
				BedType:  randBedType,
				Size:     randSize,
				Price:    getRandomPrice(randInt),
				CreateAt: primitive.NewDateTimeFromTime(time.Now()),
				UpdateAt: primitive.NewDateTimeFromTime(time.Now()),
			})
		}

		insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("New Hotel: ", insertedHotel)

		for i := range rooms {
			rooms[i].HotelID = insertedHotel.ID
		}

		updatedCount, err := roomStore.InsertRooms(ctx, &rooms, insertedHotel.ID.Hex())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Update Hotel with %d rooms \n", updatedCount)
	}
}

func main() {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client, db.DBNAME)
	roomStore = db.NewMongoRoomStore(client, db.DBNAME, hotelStore)

	seedHotels(6)
}

func init() {
	fmt.Println(" Seeding the db")
}

/** ============= Helpers ============= */

func getRandomRoomType() (types.RoomType, int) {
	randInt := rand.Intn(int(types.ClosedRoomType))
	switch randInt {
	case 1:
		return types.SingleRoomType, randInt
	case 2:
		return types.DoubleRoomType, randInt
	case 3:
		return types.TripleRoomType, randInt
	case 4:
		return types.QuadRoomType, randInt
	default:
		return types.SingleRoomType, randInt
	}
}

func getRandomSizeStirng(randInt types.RoomType) string {
	switch randInt {
	case 1:
		return types.RoomSizeSmall
	case 2:
		return types.RoomSizeNormal
	case 3:
		return types.RoomSizeNormal
	case 4:
		return types.RoomSizeKingSize
	default:
		return types.RoomSizeSmall
	}
}

func getRandomBedType(randInt types.RoomType) types.BedType {
	anotherRandInt := rand.Intn(int(types.ClosedBedType))
	switch randInt {
	default:
		return types.QueenBedType
	case types.SingleRoomType:
	case types.DoubleRoomType:
		switch anotherRandInt {
		case 1:
		case 2:
			return types.KingBedType
		case 3:
			return types.TwinBedType
		default:
			return types.QueenBedType
		}
	case types.TripleRoomType:
		switch anotherRandInt {
		case 1:
		case 2:
			return types.KingBedType
		case 3:
			return types.NormalBedType
		default:
			return types.QueenBedType
		}
	case types.QuadRoomType:
		switch anotherRandInt {
		case 3:
			return types.TwinBedType
		default:
			return types.DoubleDoubleBedType
		}
	}
	return types.QueenBedType
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
