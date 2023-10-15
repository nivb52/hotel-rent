package db

import (
	"context"
	"fmt"

	"github.com/nivb52/hotel-rent/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingColl = "bookings"

type BookingStore interface {
	InsertBooking(context.Context, *types.BookingParamsForCreate) (*types.Booking, error)
	GetBookings(context.Context, *types.BookingFilter) ([]*types.Booking, error)
	GetBookingsByRoomId(context.Context, string) ([]*types.Booking, error)
	IsRoomAvailable(context.Context, *types.BookingFilter) (bool, error)
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client, dbname string) *MongoBookingStore {
	if dbname == "" {
		dbname = DBNAME
	}

	return &MongoBookingStore{
		client: client,
		coll:   client.Database(dbname).Collection(bookingColl),
	}
}

// get bookings per room
func (s *MongoBookingStore) GetBookingsByRoomId(ctx context.Context, roomID string) ([]*types.Booking, error) {
	roomOID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"roomID": roomOID}
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var bookings []*types.Booking
	err = cur.All(ctx, &bookings)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

// build query for lookup a booking
// return example bson.M{"_id": roomOID, "fromDate": {$gte: fromDate}, "tillDate": {$lte: tillDate}}
func buildFilter(filterData *types.BookingFilter) (bson.M, error) {
	filter := bson.M{}

	if filterData.RoomID != "" {
		roomOID, err := primitive.ObjectIDFromHex(filterData.RoomID)
		if err != nil {
			return nil, err
		}

		filter["roomID"] = roomOID
	}

	if !filterData.FromDate.IsZero() {
		filter["fromDate"] = bson.M{
			"$gte": filterData.FromDate,
		}
	}

	if !filterData.TillDate.IsZero() {
		filter["tillDate"] = bson.M{
			"$lte": filterData.TillDate,
		}
	}

	fmt.Println("built filter: ", filter, " Out of: ", filterData)
	return filter, nil

}

func (s *MongoBookingStore) GetBookings(ctx context.Context, where *types.BookingFilter) ([]*types.Booking, error) {
	filter, err := buildFilter(where)
	if err != nil {
		fmt.Println("Failed to build booking filter due: ", err)
		return nil, err
	}

	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var bookings []*types.Booking
	err = cur.All(ctx, &bookings)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *MongoBookingStore) IsRoomAvailable(ctx context.Context, where *types.BookingFilter) (bool, error) {
	filter, err := buildFilter(where)
	if err != nil {
		fmt.Println("Failed to build booking filter due: ", err)
		return false, err
	}

	var reserved *types.Booking
	if err = s.coll.FindOne(ctx, filter).Decode(&reserved); err != nil {
		return false, err
	}

	if reserved.ID.IsZero() {
		return false, nil
	}

	return true, nil
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, rawData *types.BookingParamsForCreate) (*types.Booking, error) {

	RoomOID, err := primitive.ObjectIDFromHex(rawData.RoomID)
	if err != nil {
		return nil, err
	}

	UserOID, err := primitive.ObjectIDFromHex(rawData.UserID)
	if err != nil {
		return nil, err
	}

	var bookingData types.Booking
	bookingData.NumPersons = rawData.NumPersons
	bookingData.FromDate = rawData.FromDate
	bookingData.TillDate = rawData.TillDate

	bookingData.RoomID = RoomOID
	bookingData.UserID = UserOID

	res, err := s.coll.InsertOne(ctx, bookingData)
	if err != nil {
		return nil, err
	}

	bookingData.ID = res.InsertedID.(primitive.ObjectID)
	return &bookingData, nil
}
