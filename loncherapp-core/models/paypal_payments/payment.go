package paypal_payments

import "time"

type Payment struct {
	ID         string    `json:"id,omitempty" db:"ID"`
	LoncheraID int32     `json:"lonchera_id" db:"Lonchera_ID"`
	Amount     float32   `json:"amount" db:"Amount"`
	Currency   string    `json:"currency" db:"Currency"`
	Status     string    `json:"status" db:"Status"`
	CreatedAt  time.Time `json:"created_at" db:"Created_At"`
	Type       string    `json:"type" db:"Type"`
}
