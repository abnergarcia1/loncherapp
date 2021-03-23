package tracking

import "time"

type LoncheraLocation struct {
	ProfileID int
	Latitude  float32
	Longitude float32
	CreatedAt time.Time
}
