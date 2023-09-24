package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nivb52/hotel-rent/types"
)

const roomColl = "rooms"

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	InsertRooms(context.Context, *[]types.Room, string) (int, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, dbname string, HotelStore *MongoHotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(dbname).Collection(roomColl),
		HotelStore: HotelStore,
	}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		fmt.Printf("Inserting an Room Failed, reason: %s", err)
		return &types.Room{}, err
	}

	if insertedID, ok := res.InsertedID.(primitive.ObjectID); ok {
		room.ID = insertedID
		err = s.HotelStore.AddHotelRoom(ctx, room.HotelID.Hex(), room.ID.Hex())
		if err != nil {
			fmt.Printf("Failed:: res.InsertedID is not a primitive.ObjectID: %v", res.InsertedID)
			return nil, err
		}
	} else {
		fmt.Printf("Failed:: res.InsertedID is not a primitive.ObjectID: %v", res.InsertedID)
		// Handle the case where the type assertion failed (?)
	}

	return room, nil
}

// function to insert many rooms and update the relevant hotel
func (s *MongoRoomStore) InsertRooms(ctx context.Context, rooms *[]types.Room, hotelID string) (int, error) {
	hotelOID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		fmt.Printf("Failed:: supplied hotelID cannot be transformed to primitive.ObjectID: %v \n", hotelID)
		return 0, err
	}

	// https://github.com/golang/go/wiki/InterfaceSlice
	var roomsAsInterfaceSlice []interface{} = make([]interface{}, len(*rooms))
	for i, room := range *rooms {
		roomsAsInterfaceSlice[i] = room
	}

	res, err := s.coll.InsertMany(ctx, roomsAsInterfaceSlice)
	if err != nil {
		fmt.Printf("Inserting an Rooms Failed, reason: %s", err)
		return 0, err
	}

	verefiedInsertedIDs := make([]primitive.ObjectID, len(res.InsertedIDs))
	for i, oid := range res.InsertedIDs {
		if insertedID, ok := oid.(primitive.ObjectID); ok {
			verefiedInsertedIDs[i] = insertedID
		} else {
			fmt.Printf("Failed:: res.InsertedID is not a primitive.ObjectID: %v \n", res.InsertedIDs)
			// Handle the case where the type assertion failed (?)
		}
	}

	err = s.HotelStore.AddHotelRooms(ctx, hotelOID, &verefiedInsertedIDs)
	if err != nil {
		fmt.Printf("Inserting an Rooms Failed, reason: %s", err)
		return 0, err
	}

	if len(verefiedInsertedIDs) != len(*rooms) {
		err = fmt.Errorf("Some Rooms Insert Failed, wanted: %d, updated: %d", len(*rooms), len(verefiedInsertedIDs))
		return len(verefiedInsertedIDs), err
	}

	return len(verefiedInsertedIDs), nil
}
