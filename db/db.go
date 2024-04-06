package db

const (
	// uses only if db connection string is not defined in the env file
	DEFAULT_DBURI = "mongodb://root:myrootpassword@localhost:27017"
	DBNAME        = "hotel-rent"
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func GetDBUri(dburi string, customDburi ...string) string {

	if len(dburi) < 1 {
		if len(customDburi) == 1 && len(customDburi[0]) > 1 {
			dburi = customDburi[0]
		} else {
			dburi = DEFAULT_DBURI
		}
	}
	return dburi
}
