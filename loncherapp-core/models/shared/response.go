package shared

type SimpleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
