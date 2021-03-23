package menus

import "time"

type Menu struct {
	ID          int32     `json:"ID,omitempty" db:"ID"`
	LoncheraID  int32     `json:"lonchera_id,omitempty" db:"Lonchera_ID"`
	Name        string    `json:"name,omitempty" db:"Name" validate:"required"`
	Description string    `json:"description,omitempty" db:"Description"`
	Price       float64   `json:"price,omitempty" db:"Price" "`
	Currency    string    `json:"currency,omitempty" db:"Currency" validate:"max=3`
	ImageURL    string    `json:"image_url,omitempty" db:"Image_URL"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"Created_At" time_format:"sql_datetime" time_location:"UTC"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"Updated_At" time_format:"sql_datetime" time_location:"UTC"`
}
