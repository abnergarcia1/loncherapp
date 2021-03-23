package models

import "time"

type User struct {
	ID            int32     `json:"id" db:"ID"`
	TypeID        uint64    `json:"type_id" db:"Type_ID"`
	FirstName     string    `json:"first_name" db:"FirstName"`
	LastName      string    `json:"last_name" db:"LastName"`
	Email         string    `json:"email" db:"Email"`
	CreationDate  time.Time `json:"creation_date" db:"CreationDate"`
	UpdatedDate   time.Time `json:"updated_date" db:"UpdatedDate"`
	Active        bool      `json:"active" db:"Active"`
	Password      string    `json:"password,omitempty" db:"Password"`
	ProfileID     int32     `json:"profile_id,omitempty" db:"Profile_ID,omitempty"`
	ProfileActive bool      `json:"profile_active,omitempty" db:"Profile_Active,omitempty"`
}
