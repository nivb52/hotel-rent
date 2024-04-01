package mock

import (
	"math/rand"
	"time"

	"github.com/nivb52/hotel-rent/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/** ============= Data ============= */

func Hotel() types.Hotel {
	hotelInputOptions := MockHotelsMainData(52)
	int := RandomIntByMaxAndMin(51, 0)
	sub := hotelInputOptions[int : int+1]
	input := sub[0]
	hotel := MockHotelByInput(&input)
	return hotel
}

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

func User() types.UserRequiredData {
	userInputs := MockUsersMainData(40)
	int := RandomIntByMaxAndMin(40, 0)
	userDereference := *userInputs
	sub := userDereference[int : int+1]
	return sub[0]
}

func Booking(userID string, roomID string, intSeeds ...int) types.BookingParamsForCreate {
	randomDays := 5
	randomPersonsMax := 4
	randomPersonsMin := 2

	for idx, val := range intSeeds {
		if idx == 0 {
			randomDays = val
		}
		if idx == 1 {
			randomPersonsMax = val
		}
		if idx == 2 {
			randomPersonsMin = val
		}

	}

	from := time.Now()
	till := from.AddDate(0, 0, randomDays)
	bookingData := types.BookingParamsForCreate{
		UserID:     userID,
		RoomID:     roomID,
		FromDate:   from,
		TillDate:   till,
		NumPersons: RandomIntByMaxAndMin(randomPersonsMax, randomPersonsMin),
	}

	return bookingData
}

func MockUsersMainData(numberOfUsers int) *[]types.UserRequiredData {

	usersMock := []types.UserRequiredData{
		{FName: "John", LName: "Smith", Email: "john.smith@email.com"},
		{FName: "Sarah", LName: "Johnson", Email: "sarah.johnson@email.com"},
		{FName: "Michael", LName: "Brown", Email: "michael.brown@email.com"},
		{FName: "Emily", LName: "Davis", Email: "emily.davis@email.com"},
		{FName: "David", LName: "Wilson", Email: "david.wilson@email.com"},
		{FName: "Jennifer", LName: "Lee", Email: "jennifer.lee@email.com"},
		{FName: "James", LName: "Anderson", Email: "james.anderson@email.com"},
		{FName: "Jessica", LName: "Martinez", Email: "jessica.martinez@email.com"},
		{FName: "Robert", LName: "Thompson", Email: "robert.thompson@email.com"},
		{FName: "Laura", LName: "Garcia", Email: "laura.garcia@email.com"},
		{FName: "William", LName: "Taylor", Email: "william.taylor@email.com"},
		{FName: "Olivia", LName: "White", Email: "olivia.white@email.com"},
		{FName: "Benjamin", LName: "Harris", Email: "benjamin.harris@email.com"},
		{FName: "Emma", LName: "Clark", Email: "emma.clark@email.com"},
		{FName: "Christopher", LName: "Hall", Email: "christopher.hall@email.com"},
		{FName: "Sophia", LName: "Turner", Email: "sophia.turner@email.com"},
		{FName: "Matthew", LName: "Moore", Email: "matthew.moore@email.com"},
		{FName: "Ava", LName: "Parker", Email: "ava.parker@email.com"},
		{FName: "Joseph", LName: "Evans", Email: "joseph.evans@email.com"},
		{FName: "Grace", LName: "Carter", Email: "grace.carter@email.com"},
		{FName: "Samuel", LName: "Hughes", Email: "samuel.hughes@email.com"},
		{FName: "Lily", LName: "Bennett", Email: "lily.bennett@email.com"},
		{FName: "Daniel", LName: "Hayes", Email: "daniel.hayes@email.com"},
		{FName: "Mia", LName: "Foster", Email: "mia.foster@email.com"},
		{FName: "Andrew", LName: "Rivera", Email: "andrew.rivera@email.com"},
		{FName: "Harper", LName: "Gray", Email: "harper.gray@email.com"},
		{FName: "Ethan", LName: "Russell", Email: "ethan.russell@email.com"},
		{FName: "Chloe", LName: "Coleman", Email: "chloe.coleman@email.com"},
		{FName: "Christopher", LName: "Powell", Email: "christopher.powell@email.com"},
		{FName: "Isabella", LName: "Mitchell", Email: "isabella.mitchell@email.com"},
		{FName: "Alexander", LName: "Perry", Email: "alexander.perry@email.com"},
		{FName: "Madison", LName: "Ward", Email: "madison.ward@email.com"},
		{FName: "Samuel", LName: "Nelson", Email: "samuel.nelson@email.com"},
		{FName: "Elizabeth", LName: "Ramirez", Email: "elizabeth.ramirez@email.com"},
		{FName: "Nicholas", LName: "Jenkins", Email: "nicholas.jenkins@email.com"},
		{FName: "Charlotte", LName: "Wood", Email: "charlotte.wood@email.com"},
		{FName: "Ryan", LName: "Hayes", Email: "ryan.hayes@email.com"},
		{FName: "Avery", LName: "Green", Email: "avery.green@email.com"},
		{FName: "Joshua", LName: "Cooper", Email: "joshua.cooper@email.com"},
		{FName: "Sofia", LName: "Ross", Email: "sofia.ross@email.com"},
	}

	res := usersMock[0:numberOfUsers]
	return &res
}

func MockHotelsMainData(numberOfHotels int) []types.HotelRequiredData {
	hotelsMock := []types.HotelRequiredData{
		{Name: "The Ritz-Carlton", Location: "Los Angeles, California", Rating: 5},
		{Name: "Grand Hyatt", Location: "New York City, New York", Rating: 4},
		{Name: "Marriott Marquis", Location: "Atlanta, Georgia", Rating: 3},
		{Name: "Four Seasons Resort", Location: "Maui, Hawaii", Rating: 5},
		{Name: "Hilton Garden Inn", Location: "Chicago, Illinois", Rating: 3},
		{Name: "Fairmont Empress", Location: "Victoria, British Columbia, Canada", Rating: 5},
		{Name: "The Venetian", Location: "Las Vegas, Nevada", Rating: 4},
		{Name: "Burj Al Arab Jumeirah", Location: "Dubai, United Arab Emirates", Rating: 5},
		{Name: "The Savoy", Location: "London, United Kingdom", Rating: 3},
		{Name: "Hotel del Coronado", Location: "San Diego, California", Rating: 4},
		{Name: "The Ritz-Carlton", Location: "Los Angeles, California", Rating: 5},
		{Name: "Marriott Downtown", Location: "New York City, New York", Rating: 4},
		{Name: "Four Seasons Hotel", Location: "Miami, Florida", Rating: 5},
		{Name: "Hilton Riverside", Location: "New Orleans, Louisiana", Rating: 4},
		{Name: "Hyatt Regency", Location: "Chicago, Illinois", Rating: 4},
		{Name: "Sheraton Grand", Location: "Phoenix, Arizona", Rating: 4},
		{Name: "Waldorf Astoria", Location: "Beverly Hills, California", Rating: 5},
		{Name: "The Westin", Location: "Seattle, Washington", Rating: 4},
		{Name: "Fairmont", Location: "San Francisco, California", Rating: 5},
		{Name: "InterContinental", Location: "Atlanta, Georgia", Rating: 4},
		{Name: "JW Marriott", Location: "Las Vegas, Nevada", Rating: 4},
		{Name: "The Peninsula", Location: "Chicago, Illinois", Rating: 5},
		{Name: "Renaissance", Location: "Austin, Texas", Rating: 4},
		{Name: "Mandarin Oriental", Location: "Boston, Massachusetts", Rating: 5},
		{Name: "Omni", Location: "Nashville, Tennessee", Rating: 2},
		{Name: "Ritz-Carlton", Location: "Dallas, Texas", Rating: 5},
		{Name: "Hilton Garden Inn", Location: "Orlando, Florida", Rating: 3},
		{Name: "Grand Hyatt", Location: "San Diego, California", Rating: 4},
		{Name: "Kimpton", Location: "Portland, Oregon", Rating: 4},
		{Name: "Radisson Blu", Location: "Minneapolis, Minnesota", Rating: 4},
		{Name: "DoubleTree by Hilton", Location: "Philadelphia, Pennsylvania", Rating: 4},
		{Name: "Hilton", Location: "Denver, Colorado", Rating: 3},
		{Name: "Aloft", Location: "Raleigh, North Carolina", Rating: 3},
		{Name: "The Langham", Location: "New York City, New York", Rating: 1},
		{Name: "Embassy Suites", Location: "San Antonio, Texas", Rating: 3},
		{Name: "Ritz-Carlton", Location: "Naples, Florida", Rating: 5},
		{Name: "Hilton", Location: "Indianapolis, Indiana", Rating: 4},
		{Name: "The Waldorf Astoria", Location: "Las Vegas, Nevada", Rating: 5},
		{Name: "The Omni", Location: "Louisville, Kentucky", Rating: 4},
		{Name: "The St. Regis", Location: "Washington, D.C.", Rating: 5},
		{Name: "Fairfield Inn & Suites", Location: "Houston, Texas", Rating: 3},
		{Name: "Hilton", Location: "Tampa, Florida", Rating: 4},
		{Name: "Hyatt", Location: "San Jose, California", Rating: 4},
		{Name: "Marriott", Location: "St. Louis, Missouri", Rating: 1},
		{Name: "Sheraton", Location: "Salt Lake City, Utah", Rating: 4},
		{Name: "The Ritz-Carlton", Location: "Maui, Hawaii", Rating: 5},
		{Name: "InterContinental", Location: "Miami, Florida", Rating: 4},
		{Name: "JW Marriott", Location: "Orlando, Florida", Rating: 4},
		{Name: "Hilton Garden Inn", Location: "Louisville, Kentucky", Rating: 3},
		{Name: "Four Seasons", Location: "New York City, New York", Rating: 5},
		{Name: "The Westin", Location: "Denver, Colorado", Rating: 4},
		{Name: "The Peninsula", Location: "Beverly Hills, California", Rating: 5},
		{Name: "Waldorf Astoria", Location: "Orlando, Florida", Rating: 5},
		{Name: "Renaissance", Location: "Las Vegas, Nevada", Rating: 2},
		{Name: "Hyatt Regency", Location: "Seattle, Washington", Rating: 4},
		{Name: "DoubleTree by Hilton", Location: "San Diego, California", Rating: 4},
		{Name: "Kimpton", Location: "Chicago, Illinois", Rating: 4},
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
