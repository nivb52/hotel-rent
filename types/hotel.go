package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating   int8                 `bson:"rating" json:"rating"` // Rating by the rating oganization (not reviwers)

	CreateAt primitive.DateTime `bson:"create_at" json:"createAt"`
	UpdateAt primitive.DateTime `bson:"update_at" json:"updateAt"`
}

type RoomType int

const (
	_ RoomType = iota
	SingleRoomType
	DoubleRoomType
	SeaSideRoomType
	DeluxRoomType
	ClosedRoomType // Not For Reservations @attention: keep it last
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	HotelID   primitive.ObjectID `bson:"hotelID,omitempty " json:"hotelID,omitempty"`
	Type      RoomType           `bson:"type" json:"type"`
	BasePrice int                `bson:"base_price" json:"basePrice"`
	Price     int                `bson:"price" json:"price"` // we can use it as promotion price
	CreateAt  primitive.DateTime `bson:"create_at" json:"createAt"`
	UpdateAt  primitive.DateTime `bson:"update_at" json:"updateAt"`
}
