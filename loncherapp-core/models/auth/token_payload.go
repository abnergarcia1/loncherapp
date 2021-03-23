package models

type TokenPayload struct {
	UserID    string `json:"user_id"`
	ProfileID string `json:"profile_id"`
	UUID      string `json:"-"`
}
