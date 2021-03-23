package routes

import "time"

type Schedule struct {
	ID        int32     `json:"id,omitempty" db:"ID"`
	RouteID   int32     `json:"route_id" db:"Route_ID"`
	Weekday   int32     `json:"weekday" db:"Weekday"`
	ArriveAt  time.Time `json:"arrive_at" db:"Arrive_At" time_format:"sql_datetime" time_location:"UTC"`
	GoneAt    time.Time `json:"gone_at" db:"Gone_At" time_format:"sql_datetime" time_location:"UTC"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"Created_At" time_format:"sql_datetime" time_location:"UTC"`
	Active    bool      `json:"active,omitempty" db:"Active"`
	ProfileID int32     `json:"-"`
}
