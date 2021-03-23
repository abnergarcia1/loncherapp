package shared

type AppConfigParam struct {
	ID        int32  `json:"id,omitempty" db:"ID,omitempty"`
	Parameter string `json:"parameter" db:"Parameter"`
	Value     string `json:"value" db:"Value"`
}
