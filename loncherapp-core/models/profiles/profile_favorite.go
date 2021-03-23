package profiles

import "time"

type Favorite struct {
	UserID    int32     `json:"user_id,omitempty" db:"User_ID"`
	ProfileID int32     `json:"profile_id" db:"Lonchera_ID"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"Created_At"`
}
