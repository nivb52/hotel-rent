package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nivb52/hotel-rent/types"
)

const roomColl = "rooms"

type RoomsFilter struct {
	BasePrice int
}

type RoomStore interface {
	GetRoomByIDs(context.Context, *[]string) (*[]types.Room, error)
	//  2nd Argument - hotel ID
	GetHotelRooms(context.Context, string, *types.HotelFilter) (*[]types.Room, error)
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	InsertRooms(context.Context, *[]types.Room, string) (int, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, HotelStore *MongoHotelStore, dbname string) *MongoRoomStore {
	if dbname == "" {
		dbname = DBNAME
	}
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(dbname).Collection(roomColl),
		HotelStore: HotelStore,
	}
}

// function to retrive rooms data using the hotel ID
func (s *MongoRoomStore) GetHotelRooms(ctx context.Context, id string, where *types.HotelFilter) (*[]types.Room, error) {

	filter := buildHotelRoomsFilter(where)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter["hotelID"] = oid

	var rooms []types.Room
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &rooms)
	if err != nil {
		return nil, err
	}

	return &rooms, nil
}

func (s *MongoRoomStore) GetRoomByIDs(ctx context.Context, roomIDs *[]string) (*[]types.Room, error) {
	var roomOIDs []primitive.ObjectID = make([]primitive.ObjectID, len(*roomIDs))
	var errors []error = make([]error, len(*roomIDs))

	for i, id := range *roomIDs {
		oid, err := primitive.ObjectIDFromHex(id)
		roomOIDs[i] = oid
		if err != nil {
			errors = append(errors, err)
		}
	}

	var rooms []types.Room
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &rooms)
	if err != nil {
		return nil, err
	}

	if len(errors) > 0 {
		return &rooms, errors[0]
	}

	return &rooms, nil
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

func buildHotelRoomsFilter(filterData *types.HotelFilter) bson.M {
	filter := bson.M{}

	if filterData.MaxPrice > 0 && filterData.MinPrice > 0 {
		filter["price"] = bson.M{
			"$gte": filterData.MaxPrice,
			"$lte": filterData.MaxPrice,
		}
	} else {
		if filterData.MaxPrice > 0 && filterData.MinPrice == 0 {
			filter["price"] = bson.M{
				"$gte": filterData.MaxPrice,
			}
		} else if filterData.MinPrice > 0 && filterData.MaxPrice == 0 {
			filter["price"] = bson.M{
				"$lte": filterData.MinPrice,
			}
		}
	}

	if filterData.BedType > 0 && filterData.BedType < int(types.ClosedBedType) {
		filter["bedType"] = filterData.BedType
	}

	if len(filterData.RoomSize) > 0 &&
		(filterData.RoomSize == types.RoomSizeKingSize ||
			filterData.RoomSize == types.RoomSizeNormal ||
			filterData.RoomSize == types.RoomSizeSmall) {
		filter["size"] = filterData.RoomSize
	}

	return filter
}
