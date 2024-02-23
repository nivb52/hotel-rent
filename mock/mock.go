package mock

import (
	"math/rand"
	"time"

	"github.com/nivb52/hotel-rent/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/** ============= Data ============= */

func MockHotelByInput(input *types.HotelRequiredData) types.Hotel {
	hotel := types.Hotel{
		Name:     input.Name,
		Location: input.Location,
		Rating:   input.Rating,
		CreateAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdateAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	return hotel
}

func MockRoom(numberOfRooms int) []types.Room {
	rooms := []types.Room{}
	if numberOfRooms < 1 {
		numberOfRooms = 1
	}

	for i := 0; i < numberOfRooms; i++ {
		randRoomType, randInt := getRandomRoomType()
		randBedType := getRandomBedType(types.RoomType(randInt))
		randSize := getRandomSizeString(types.RoomType(randInt))

		rooms = append(rooms, types.Room{
			Type:     randRoomType,
			BedType:  randBedType,
			Size:     randSize,
			Price:    getRandomPrice(randInt),
			CreateAt: primitive.NewDateTimeFromTime(time.Now()),
			UpdateAt: primitive.NewDateTimeFromTime(time.Now()),
		})
	}

	return rooms
}

func MockUsersMainData(numberOfUsers int) *[]types.UserRequiredData {

	usersMock := []types.UserRequiredData{
		{"John", "Smith", "john.smith@email.com"},
		{"Sarah", "Johnson", "sarah.johnson@email.com"},
		{"Michael", "Brown", "michael.brown@email.com"},
		{"Emily", "Davis", "emily.davis@email.com"},
		{"David", "Wilson", "david.wilson@email.com"},
		{"Jennifer", "Lee", "jennifer.lee@email.com"},
		{"James", "Anderson", "james.anderson@email.com"},
		{"Jessica", "Martinez", "jessica.martinez@email.com"},
		{"Robert", "Thompson", "robert.thompson@email.com"},
		{"Laura", "Garcia", "laura.garcia@email.com"},
		{"William", "Taylor", "william.taylor@email.com"},
		{"Olivia", "White", "olivia.white@email.com"},
		{"Benjamin", "Harris", "benjamin.harris@email.com"},
		{"Emma", "Clark", "emma.clark@email.com"},
		{"Christopher", "Hall", "christopher.hall@email.com"},
		{"Sophia", "Turner", "sophia.turner@email.com"},
		{"Matthew", "Moore", "matthew.moore@email.com"},
		{"Ava", "Parker", "ava.parker@email.com"},
		{"Joseph", "Evans", "joseph.evans@email.com"},
		{"Grace", "Carter", "grace.carter@email.com"},
		{"Samuel", "Hughes", "samuel.hughes@email.com"},
		{"Lily", "Bennett", "lily.bennett@email.com"},
		{"Daniel", "Hayes", "daniel.hayes@email.com"},
		{"Mia", "Foster", "mia.foster@email.com"},
		{"Andrew", "Rivera", "andrew.rivera@email.com"},
		{"Harper", "Gray", "harper.gray@email.com"},
		{"Ethan", "Russell", "ethan.russell@email.com"},
		{"Chloe", "Coleman", "chloe.coleman@email.com"},
		{"Christopher", "Powell", "christopher.powell@email.com"},
		{"Isabella", "Mitchell", "isabella.mitchell@email.com"},
		{"Alexander", "Perry", "alexander.perry@email.com"},
		{"Madison", "Ward", "madison.ward@email.com"},
		{"Samuel", "Nelson", "samuel.nelson@email.com"},
		{"Elizabeth", "Ramirez", "elizabeth.ramirez@email.com"},
		{"Nicholas", "Jenkins", "nicholas.jenkins@email.com"},
		{"Charlotte", "Wood", "charlotte.wood@email.com"},
		{"Ryan", "Hayes", "ryan.hayes@email.com"},
		{"Avery", "Green", "avery.green@email.com"},
		{"Joshua", "Cooper", "joshua.cooper@email.com"},
		{"Sofia", "Ross", "sofia.ross@email.com"},
	}

	res := usersMock[0:numberOfUsers]
	return &res
}

func MockHotelsMainData(numberOfHotels int) []types.HotelRequiredData {
	hotelsMock := []types.HotelRequiredData{
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
		{"The Ritz-Carlton", "Los Angeles, California", 5},
		{"Marriott Downtown", "New York City, New York", 4},
		{"Four Seasons Hotel", "Miami, Florida", 5},
		{"Hilton Riverside", "New Orleans, Louisiana", 4},
		{"Hyatt Regency", "Chicago, Illinois", 4},
		{"Sheraton Grand", "Phoenix, Arizona", 4},
		{"Waldorf Astoria", "Beverly Hills, California", 5},
		{"The Westin", "Seattle, Washington", 4},
		{"Fairmont", "San Francisco, California", 5},
		{"InterContinental", "Atlanta, Georgia", 4},
		{"JW Marriott", "Las Vegas, Nevada", 4},
		{"The Peninsula", "Chicago, Illinois", 5},
		{"Renaissance", "Austin, Texas", 4},
		{"Mandarin Oriental", "Boston, Massachusetts", 5},
		{"Omni", "Nashville, Tennessee", 2},
		{"Ritz-Carlton", "Dallas, Texas", 5},
		{"Hilton Garden Inn", "Orlando, Florida", 3},
		{"Grand Hyatt", "San Diego, California", 4},
		{"Kimpton", "Portland, Oregon", 4},
		{"Radisson Blu", "Minneapolis, Minnesota", 4},
		{"DoubleTree by Hilton", "Philadelphia, Pennsylvania", 4},
		{"Hilton", "Denver, Colorado", 3},
		{"Aloft", "Raleigh, North Carolina", 3},
		{"The Langham", "New York City, New York", 1},
		{"Embassy Suites", "San Antonio, Texas", 3},
		{"Ritz-Carlton", "Naples, Florida", 5},
		{"Hilton", "Indianapolis, Indiana", 4},
		{"The Waldorf Astoria", "Las Vegas, Nevada", 5},
		{"The Omni", "Louisville, Kentucky", 4},
		{"The St. Regis", "Washington, D.C.", 5},
		{"Fairfield Inn & Suites", "Houston, Texas", 3},
		{"Hilton", "Tampa, Florida", 4},
		{"Hyatt", "San Jose, California", 4},
		{"Marriott", "St. Louis, Missouri", 1},
		{"Sheraton", "Salt Lake City, Utah", 4},
		{"The Ritz-Carlton", "Maui, Hawaii", 5},
		{"InterContinental", "Miami, Florida", 4},
		{"JW Marriott", "Orlando, Florida", 4},
		{"Hilton Garden Inn", "Louisville, Kentucky", 3},
		{"Four Seasons", "New York City, New York", 5},
		{"The Westin", "Denver, Colorado", 4},
		{"The Peninsula", "Beverly Hills, California", 5},
		{"Waldorf Astoria", "Orlando, Florida", 5},
		{"Renaissance", "Las Vegas, Nevada", 2},
		{"Hyatt Regency", "Seattle, Washington", 4},
		{"DoubleTree by Hilton", "San Diego, California", 4},
		{"Kimpton", "Chicago, Illinois", 4},
	}

	res := hotelsMock[0:numberOfHotels]
	return res
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

func getRandomSizeString(randRoomType types.RoomType) string {
	switch randRoomType {
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

func getRandomBedType(randRoomType types.RoomType) types.BedType {
	bedTypeRand := rand.Intn(int(types.ClosedBedType))
	switch randRoomType {
	default:
		return types.QueenBedType
	case types.SingleRoomType:
	case types.DoubleRoomType:
		switch bedTypeRand {
		case 1:
		case 2:
			return types.KingBedType
		case 3:
			return types.TwinBedType
		default:
			return types.QueenBedType
		}
	case types.TripleRoomType:
		switch bedTypeRand {
		case 1:
		case 2:
			return types.KingBedType
		case 3:
			return types.NormalBedType
		default:
			return types.QueenBedType
		}
	case types.QuadRoomType:
		switch bedTypeRand {
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
		return RandomIntByMaxAndMin(99, 50) //types.SingleRoomType
	case 2:
		return RandomIntByMaxAndMin(150, 100) //types.DoubleRoomType
	case 3:
		return RandomIntByMaxAndMin(200, 150) //types.SeaSideRoomType
	case 4:
		return RandomIntByMaxAndMin(500, 200) //types.DeluxRoomType
	default:
		return RandomIntByMaxAndMin(100, 50) //types.SingleRoomType
	}
}

// function which return random integer between to max and min
func RandomIntByMaxAndMin(max int, min int) int {
	return rand.Intn(max-min) + min
}
