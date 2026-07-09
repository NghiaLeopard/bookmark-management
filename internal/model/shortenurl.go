package model

// ShortenUrlResponse is returned when a short code is created successfully.
type ShortenUrlResponse struct {
	Code    string `json:"code" example:"a1b2c3d"`
	Message string `json:"message" example:"Shorten URL generated successfully!"`
}
