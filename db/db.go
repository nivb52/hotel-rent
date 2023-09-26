package db

const (
	DBURI      = "mongodb://localhost:27017"
	DBNAME     = "hotel-rent"
	TestDBNAME = "hotel-rent-testing"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
