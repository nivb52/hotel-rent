package db

import (
	"context"

	"github.com/nivb52/hotel-rent/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingColl = "bookings"

type BookingStore interface {
	InsertBooking(context.Context, *types.BookingParamsForCreate) (*types.Booking, error)
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

func (s *MongoBookingStore) InsertBooking(ctx context.Context, data *types.BookingParamsForCreate) (*types.Booking, error) {

	res, err := s.coll.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	var newBooking types.Booking
	newBooking.NumPersons = data.NumPersons
	newBooking.FromDate = data.FromDate
	newBooking.TillDate = data.TillDate
	newBooking.ID = res.InsertedID.(primitive.ObjectID)

	RoomOID, err := primitive.ObjectIDFromHex(data.RoomID)
	if err != nil {
		return nil, err
	}

	UserOID, err := primitive.ObjectIDFromHex(data.UserID)
	if err != nil {
		return nil, err
	}

	newBooking.RoomID = RoomOID
	newBooking.UserID = UserOID
	return &newBooking, nil
}
