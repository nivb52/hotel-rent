package db

import (
	"context"
	"errors"
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
	GetBookingsById(ctx context.Context, ID string) (*types.Booking, error)
	GetBookingsByRoomId(context.Context, string) ([]*types.Booking, error)
	IsRoomAvailable(context.Context, *types.BookingFilter) (bool, error)
	CancelBooking(context.Context, string) error
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

// get booking by id
func (s *MongoBookingStore) GetBookingsById(ctx context.Context, ID string) (*types.Booking, error) {
	isExistID := primitive.NilObjectID.IsZero()
	if !isExistID {
		return nil, errors.New("missing Booking ID")
	}

	OID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	var booking types.Booking
	err = s.coll.FindOne(ctx, bson.M{"_id": OID}).Decode(&booking)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &booking, nil
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

func (s *MongoBookingStore) GetBookings(ctx context.Context, where *types.BookingFilter) ([]*types.Booking, error) {
	filter, err := buildBookingFilter(where)
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
	filter, err := buildBookingFilter(where)
	if err != nil {
		fmt.Println("Failed to build booking filter due: ", err)
		return false, err
	}

	var reserved *types.Booking
	err = s.coll.FindOne(ctx, filter).Decode(&reserved)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return true, nil
		}

		return false, err
	}

	if reserved.ID.IsZero() {
		return true, nil
	}

	return false, nil
}

// ## Updates
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

func (s *MongoBookingStore) CancelBooking(ctx context.Context, id string) error {
	bookOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	values := bson.M{"isCanceled": true}
	_, err = s.coll.UpdateOne(ctx,
		bson.M{"_id": bookOID}, // filter
		bson.M{"$set": values},
	)

	if err != nil {
		fmt.Println("Failed to delete booking, due: ", err)
		return err
	}

	return nil
}

// ## Helpers

// build query for lookup a booking, includes filter by dates to find existing booking
func buildBookingFilter(filterData *types.BookingFilter) (bson.M, error) {
	filter := bson.M{}

	if filterData.RoomID != "" {
		roomOID, err := primitive.ObjectIDFromHex(filterData.RoomID)
		if err != nil {
			return nil, err
		}

		filter["roomID"] = roomOID
	}

	// Find Apartment which ALREADY reserved:
	// fromDate       	   tillDate    start----end new reservation
	//  <------------------------->    *--------**     <-- should be allowed to book
	//  start---end     fromDate  		   tillDate
	//  *--------** <----------------------------->    <-- should be allowed to book
	// fromDate        new reservation    tillDate
	//  <----------------*-----------**------------>   <-- should NOT be allowed to book
	//  start    fromDate   end reservation   tillDate
	//  *--------------<------**------------------>   <-- should NOT be allowed to book
	//  fromDate   start reservation   tillDate   end reservation
	//  <----------------*--------------------->----**   <-- should NOT be allowed to book
	if !filterData.FromDate.IsZero() && !filterData.TillDate.IsZero() {
		if filter["$or"] == nil {
			filter["$or"] = []bson.M{}
		}
		startDate := (filterData.FromDate)
		endDate := (filterData.TillDate)

		filter["$or"] = append(filter["$or"].([]bson.M),
			bson.M{"$and": []bson.M{
				{"fromDate": bson.M{"$gte": startDate}},
				{"fromDate": bson.M{"$lte": endDate}},
			}},
			bson.M{"$and": []bson.M{
				{"tillDate": bson.M{"$gte": startDate}},
				{"tillDate": bson.M{"$lte": endDate}},
			}},
			bson.M{"$and": []bson.M{
				{"fromDate": bson.M{"$lte": startDate}},
				{"tillDate": bson.M{"$gte": endDate}},
			}},
		)
	} else if !filterData.FromDate.IsZero() {
		if filter["$and"] == nil {
			filter["$and"] = []bson.M{}
		}
		endDate := (filterData.TillDate)
		filter["$and"] = append(filter["$and"].([]bson.M),
			bson.M{"tillDate": bson.M{"$lte": endDate}},
		)
	} else if !filterData.TillDate.IsZero() {
		if filter["$and"] == nil {
			filter["$and"] = []bson.M{}
		}
		startDate := (filterData.FromDate)
		filter["$and"] = append(filter["$and"].([]bson.M),
			bson.M{"fromDate": bson.M{"$gte": startDate}},
		)
	}

	fmt.Println("built filter: ", filter, " Out of: ", filterData)
	return filter, nil
}

// build query for lookup a booking
func buildFreeDatesBookingFilter(filterData *types.BookingFilter) (bson.M, error) {
	filter := bson.M{}

	if filterData.RoomID != "" {
		roomOID, err := primitive.ObjectIDFromHex(filterData.RoomID)
		if err != nil {
			return nil, err
		}

		filter["roomID"] = roomOID
	}

	if filter["$and"] == nil {
		filter["$and"] = []bson.M{}
	}
	// Find Apartment which IS NOT reserved:
	// fromDate       	   tillDate    start----end new reservation
	//  <------------------------->    *--------**     <-- should be allowed to book
	//  start---end     fromDate  		   tillDate
	//  *--------** <----------------------------->    <-- should be allowed to book
	// fromDate        new reservation    tillDate
	//  <----------------*-----------**------------>   <-- should NOT be allowed to book
	//  start    fromDate   end reservation   tillDate
	//  *--------------<------**------------------>   <-- should NOT be allowed to book
	//  fromDate   start reservation   tillDate   end reservation
	//  <----------------*--------------------->----**   <-- should NOT be allowed to book
	if !filterData.FromDate.IsZero() && !filterData.TillDate.IsZero() {
		startDate := (filterData.FromDate)
		endDate := (filterData.TillDate)
		filter["$and"] = append(filter["$and"].([]bson.M),
			bson.M{"$or": []bson.M{
				{"fromDate": bson.M{"$gt": endDate}},
			}},
			bson.M{"$or": []bson.M{
				{"tillDate": bson.M{"$lt": startDate}},
			}},
		)
	} else if !filterData.FromDate.IsZero() {
		startDate := (filterData.FromDate)
		filter["$and"] = append(filter["$and"].([]bson.M),
			bson.M{"tillDate": bson.M{"$lt": startDate}},
		)
	} else if !filterData.TillDate.IsZero() {
		endDate := (filterData.TillDate)
		filter["$and"] = append(filter["$and"].([]bson.M),
			bson.M{"fromDate": bson.M{"$gt": endDate}},
		)
	}

	fmt.Println("built filter: ", filter, " Out of: ", filterData)
	return filter, nil
}
