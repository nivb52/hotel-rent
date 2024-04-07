package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelRequiredData struct {
	Name     string
	Location string
	Rating   int8
}

type Hotel struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string               `bson:"name" json:"name"`
	Location  string               `bson:"location" json:"location"`
	Rooms     []primitive.ObjectID `bson:"rooms" json:"rooms"`
	RoomsInfo []Room               `bson:"rooms_info" json:"rooms_info,omitempty"`
	Rating    int8                 `bson:"rating" json:"rating"` // Rating by the rating oganization (not reviwers)
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
	Size    string             `bson:"size" json:"size"`

	Price    int                `bson:"price" json:"price"` // we can use it as promotion price
	CreateAt primitive.DateTime `bson:"create_at" json:"createAt"`
	UpdateAt primitive.DateTime `bson:"update_at" json:"updateAt"`
}

type HotelFilter struct {
	Rating int8 `bson:"rating,omitempty" json:"rating,omitempty"`
	RoomsFilter
}

func (hf HotelFilter) HasHotelilter() bool {
	if hf.Rating > 0 {
		return true // Struct is considerd not empty
	}

	return false
}

type RoomsFilter struct {
	Rooms bool `json:"rooms,omitempty"`
	// RoomType int       `json:"type,omitempty"`
	BedType  int    `json:"bedType,omitempty"`
	RoomSize string `json:"size,omitempty"`
	MinPrice int    `json:"minPrice,omitempty"`
	MaxPrice int    `json:"maxPrice,omitempty"`
	DatesFilter
}

func (rf RoomsFilter) HasRoomsFilter() bool {
	if rf.Rooms {
		return true // Struct is considerd not empty
	}
	return false
}
