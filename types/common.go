package types

import "time"

type Map map[string]any

type DatesFilter struct {
	FromDate time.Time `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	TillDate time.Time `bson:"tillDate,omitempty" json:"tillDate,omitempty"`
}

type ResourceResp struct {
	Data  any   `json:"data"`
	Total int64 `json:"total,omitempty"`
	Page  int   `json:"page,omitempty"`
	// Total int32
}
