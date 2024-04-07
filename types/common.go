package types

import "time"

type Map map[string]any

type DatesFilter struct {
	FromDate time.Time `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	TillDate time.Time `bson:"tillDate,omitempty" json:"tillDate,omitempty"`
}
