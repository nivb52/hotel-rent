package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nivb52/hotel-rent/types"
)

const hotelColl = "hotels"

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbname string) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbname).Collection(hotelColl),
	}
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		fmt.Printf("Inserting an Hotel Failed, reason: %s", err)
		return &types.Hotel{}, err
	}

	if insertedID, ok := res.InsertedID.(primitive.ObjectID); ok {
		hotel.ID = insertedID
	} else {
		fmt.Printf("Failed:: res.InsertedID is not a primitive.ObjectID: %v", res.InsertedID)
		// Handle the case where the type assertion failed (?)
	}

	return hotel, nil

}