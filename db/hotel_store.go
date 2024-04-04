package db

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nivb52/hotel-rent/types"
)

const hotelColl = "hotels"

//@TODO: consider not use in the interface the filter as bson.M
// maybe try use []*FilterString, []*FilterInt or json ?

type HotelStore interface {
	GetHotelByID(context.Context, string) (*types.Hotel, error)
	GetHotels(context.Context) (*[]types.Hotel, error)
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotelByID(context.Context, string, *types.Hotel) (*types.Hotel, error)

	UpdateHotel(context.Context, bson.M, bson.M) error
	AddHotelRoom(context.Context, string, string) error
	AddHotelRooms(context.Context, primitive.ObjectID, *[]primitive.ObjectID) error
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbname string) *MongoHotelStore {
	if dbname == "" {
		dbname = DBNAME
	}
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbname).Collection(hotelColl),
	}
}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var hotel types.Hotel
	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel)
	if err != nil {
		return nil, err
	}

	return &hotel, nil
}

func (s *MongoHotelStore) GetHotels(ctx context.Context) (*[]types.Hotel, error) {
	var hotels []types.Hotel
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &hotels)
	if err != nil {
		return nil, err
	}

	return &hotels, nil
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

// deprecated
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

// function tp update hotel by Id
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
func (s *MongoHotelStore) AddHotelRoom(ctx context.Context, hotelID string, roomID string) error {
	hotelOid, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return err
	}

	roomOid, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return err
	}

	_, err = s.coll.UpdateOne(ctx,
		bson.M{"_id": hotelOid},
		bson.M{"$push": bson.M{"rooms": roomOid}},
	)

	if err != nil {
		errorNilInRoomField := `The field 'rooms' must be an array but is of type null in document`
		isErrorCanBeAutoHandled := strings.Contains(err.Error(), errorNilInRoomField)
		if isErrorCanBeAutoHandled == true {
			_, err = s.coll.UpdateOne(ctx,
				bson.M{"_id": hotelOid},
				bson.M{"$set": bson.M{"rooms": []primitive.ObjectID{roomOid}}},
			)
			if err != nil {
				fmt.Printf("(1st try out of 2) Failed to update room id %s in Hotel, reason: %s", roomID, err)
				fmt.Printf("(2nd try out of 2) Failed to create room array of room id %s in Hotel, reason: %s", roomID, err)
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

// function to add many rooms in hotel
func (s *MongoHotelStore) AddHotelRooms(ctx context.Context, hotelOID primitive.ObjectID, roomOIDs *[]primitive.ObjectID) error {
	_, err := s.coll.UpdateOne(ctx,
		bson.M{"_id": hotelOID},
		bson.M{"$set": bson.M{"rooms": roomOIDs}},
	)
	if err != nil {
		fmt.Printf("Failed to create room array of room in Hotel Id %s , reason: %s", hotelOID, err)
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
