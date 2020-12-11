package models

//List is the API structure for list object
type List struct {
	TotalCount int64       `json:"total_count"`
	Object     string      `json:"object"`
	Data       interface{} `json:"data"`
}
