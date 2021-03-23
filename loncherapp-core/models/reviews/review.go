package reviews

import "time"

type Review struct {
	ID         int32     `json:"id,omitempty" db:"ID"`
	LoncheraID int32     `json:"lonchera_id" db:"Lonchera_ID"`
	Comment    string    `json:"comment,omitempty" db:"Comment"`
	UserID     int32     `json:"user_id,omitempty" db:"User_ID"`
	UserName   string    `json:"user_name,omitempty" db:"User_Name"`
	Rating     int32     `json:"rating" db:"Rating" validate:"gte=0,lte=130"`
	CreatedAt  time.Time `json:"created_at,omitempty" db:"Created_At"`
}

type RatingAverage struct {
	LoncheraID int32   `json:"lonchera_id" db:"Lonchera_ID"`
	Rating     float64 `json:"rating" db:"Rating"`
}
