package models

//APIBusinessError is the Business Error structure
type APIBusinessError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Field   string `json:"field"`
}

//Error is the Error structure for all API errors
type Error struct {
	Err APIBusinessError `json:"error"`
}
