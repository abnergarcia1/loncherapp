package routes

import "time"

type Route struct {
	ID            int32      `json:"id,omitempty" db:"ID"`
	LoncheraID    int32      `json:"lonchera_id" db:"Lonchera_ID"`
	Location      string     `json:"location" db:"Location"`
	Address       string     `json:"address" db:"Address"`
	Name          string     `json:"name" db:"Name"`
	Description   string     `json:"description" db:"Description"`
	Order         int32      `json:"order" db:"Order"`
	CreatedAt     time.Time  `json:"created_at,omitempty" db:"Created_At"`
	UpdatedAt     time.Time  `json:"updated_at,omitempty" db:"Updated_At"`
	Latitude      float64    `json:"latitude" db:"Latitude"`
	Longitude     float64    `json:"longitude" db:"Longitude"`
	GooglePlaceID string     `json:"google_place_id" db:"Google_Place_ID"`
	Schedules     []Schedule `json:"schedules,omitempty" db:"-"`
}
