package types

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"userID,omitempty" json:"userID"`
	RoomID primitive.ObjectID `bson:"roomID,omitempty" json:"roomID"`

	FromDate time.Time `bson:"fromDate,omitempty" json:"fromDate"`
	TillDate time.Time `bson:"tillDate,omitempty" json:"tillDate"`
	// Number Of Guests Adults And young
	NumPersons int `bson:"numPersons" json:"numPersons"`

	// is Cofirmed, is Canceled, is Payed, privateNotes, ....
}

type BookingParamsForCreate struct {
	UserID string `bson:"userID,omitempty" json:"userID"`
	RoomID string `bson:"roomID,omitempty" json:"roomID"`

	FromDate   time.Time `bson:"fromDate,omitempty" json:"fromDate"`
	TillDate   time.Time `bson:"tillDate,omitempty" json:"tillDate"`
	NumPersons int       `bson:"numPersons" json:"numPersons"`
}

func (p BookingParamsForCreate) Validate() map[string]string {
	errors := map[string]string{}
	_, _, day := time.Now().Date()
	now := time.Now()

	// DATES
	if p.FromDate.Before(now) && p.FromDate.Day() < day {
		formatedDate := strings.Split(p.FromDate.Format(time.RFC3339), "T")[0]
		errors["fromDate"] = fmt.Sprintf("Cannot book room in the past, reviced: %s", formatedDate)
	}

	if p.FromDate.After(p.TillDate) { // also prevent TillDate < now
		formatedDate := strings.Split(p.TillDate.Format(time.RFC3339), "T")[0]
		errors["tillDate"] = fmt.Sprintf("Till date (end date) can't be before From date, reviced: %s", formatedDate)
	} else if p.FromDate.Equal(p.TillDate) {
		formatedFromDate := strings.Split(p.FromDate.Format(time.RFC3339), "T")[0]
		formatedTillDate := strings.Split(p.TillDate.Format(time.RFC3339), "T")[0]
		fmt.Printf("form: %s  til: %s", formatedFromDate, formatedTillDate)
		errors["tillDate"] = "End date can't be equal to start date"
	}

	// GUESTS
	if p.NumPersons < 1 {
		errors["NumPersons"] = fmt.Sprintf("Number of guests cannot be less then 1, reviced: %d", p.NumPersons)
	}

	if p.NumPersons > 100 {
		errors["NumPersons"] = fmt.Sprintf("Number of guests cannot be greater then 100, reviced: %d", p.NumPersons)
	}

	// Missing or Invalid Data
	if p.RoomID == "" {
		errors["RoomID"] = "Booking Data missing RoomID"
	}

	if p.UserID == "" {
		errors["UserID"] = "Booking Data missing UserID"
	}

	return errors
}
