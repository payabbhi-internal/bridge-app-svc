package models

import "net/http"

//Route is the API structure for route
type Route struct {
	Name        string
	Methods     []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
