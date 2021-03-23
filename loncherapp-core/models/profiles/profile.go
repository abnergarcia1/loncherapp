package profiles

import (
	"time"
)

type Profile struct {
	ID                int32     `json:"id,omitempty" db:"ID"`
	UserID            int32     `json:"user_id" validate:"required" db:"User_ID"`
	Description       string    `json:"description" validate:"required" db:"Description"`
	CategoryID        int32     `json:"category_id" validate:"required" db:"Category_ID"`
	CoverImageURL     string    `json:"cover_image_url" db:"Cover_Image_URL"`
	Website           string    `json:"website" db:"Website"`
	Active            bool      `json:"active,omitempty" db:"Active"`
	MembershipDueDate time.Time `json:"membership_due_date,omitempty" db:"Membership_Due_Date,omitempty"`
	CreatedAt         time.Time `json:"created_a,omitemptyt" db:"Created_At"`
	UpdatedAt         time.Time `json:"updated_at,omitempty" db:"Updated_At,omitempty"`
	//this properties belongs to another db
	Rating     float64 `json:"rating,omitempty" db:"-"`
	IsFavorite bool    `json:"is_favorite,omitempty" db:"-"`
}
