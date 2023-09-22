package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nivb52/hotel-rent/types"
)

const hotelColl = "hotels"

type FilterString struct {
	Key   string
	Value string
}

type FilterInt struct {
	Key   string
	Value int
}
type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotelByID(context.Context, string, *types.Hotel) (*types.Hotel, error)

	// TODO: instead of filter as bson.M use []*FilterString, []*FilterInt or json ?
	UpdateHotel(context.Context, bson.M, bson.M) error
	AddHotelRoom(context.Context, string, string) error
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

// function to update
func (s *MongoHotelStore) UpdateHotel(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx,
		filter,
		update,
	)

	if err != nil {
		return err
	}
	return nil
}

// update hotel by Id (will not work properly for rooms)
func (s *MongoHotelStore) UpdateHotelByID(ctx context.Context, id string, values *types.Hotel) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update, err := convertToMongoUpdateValues(values)
	if err != nil {
		fmt.Printf("Failed:: convert to bson.Map type failed: %v", err)
		return nil, err
	}

	_, err = s.coll.UpdateOne(ctx,
		bson.M{"_id": oid}, // filter
		bson.D{{Key: "$set", Value: update}},
	)

	if err != nil {
		return nil, err
	}

	return values, err
}

// function to add room in hotel
func (s *MongoHotelStore) AddHotelRoom(ctx context.Context, id string, roomId string) error {
	hotelOid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	roomOid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.coll.UpdateOne(ctx,
		bson.M{"_id": hotelOid},
		bson.M{"$push": bson.M{"rooms": roomOid}},
	)

	if err != nil {
		fmt.Printf("Failed to update room id %s in Hotel, reason: %s", roomId, err)
		return err
	}

	return nil
}

/** ============= Helpers ============= */
func convertToMongoUpdateValues(values *types.Hotel) (*bson.M, error) {
	//@note this can be improve https://stackoverflow.com/questions/66493924/converting-a-struct-to-a-bson-document
	bm, err := bson.Marshal(values)
	if err != nil {
		return nil, err
	}
	var update bson.M
	err = bson.Unmarshal(bm, &update)
	if err != nil {
		return nil, err
	}
	return &update, nil
}
