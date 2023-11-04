package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating   int8                 `bson:"rating" json:"rating"` // Rating by the rating oganization (not reviwers)
	// Reviewers   int32                `bson:"reviewers" json:"reviewers"`
	// ReviersScore

	CreateAt primitive.DateTime `bson:"create_at" json:"createAt"`
	UpdateAt primitive.DateTime `bson:"update_at" json:"updateAt"`
}

// romms and bed types
// https://hoteltechreport.com/news/room-type
// https://www.hospitality-school.com/hotel-room-types-classification/hotel-room-types-classification-2/

type RoomType int

const (
	_ RoomType = iota
	SingleRoomType
	DoubleRoomType
	TripleRoomType
	QuadRoomType
	ClosedRoomType // Not For Reservations @attention: keep it last while using seed script
)

type BedType int

const (
	_ BedType = iota
	NormalBedType
	QueenBedType
	KingBedType
	TwinBedType
	DoubleDoubleBedType
	ClosedBedType // Not For Reservations @attention: keep it last while using seed script
)

type RoomSize int

const (
	RoomSizeSmall    = "small"
	RoomSizeNormal   = "normal"
	RoomSizeKingSize = "king"
)

type Room struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	HotelID primitive.ObjectID `bson:"hotelID,omitempty " json:"hotelID,omitempty"`
	Type    RoomType           `bson:"type" json:"type"`
	BedType BedType            `bson:"bedType" json:"bedType"`
	Size    RoomSize           `bson:"size" json:"size"`

	Price    int                `bson:"price" json:"price"` // we can use it as promotion price
	CreateAt primitive.DateTime `bson:"create_at" json:"createAt"`
	UpdateAt primitive.DateTime `bson:"update_at" json:"updateAt"`
}
