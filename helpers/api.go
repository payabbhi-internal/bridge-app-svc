package helpers

import (
	"fmt"

	"github.com/paypermint/bridge-app-svc/models"
)

const (
	listObject = "list"
)

// CreateAPIResponse returns APIResponse for apilist call
func CreateAPIResponse(route models.Route) models.RouteResponse {
	return models.RouteResponse{
		Name:    route.Name,
		Methods: route.Methods,
		URL:     fmt.Sprintf("api/v1%s", route.Pattern),
	}
}

// CreateAPIList returns list of APIResponse
func CreateAPIList(routeList []models.Route) []models.RouteResponse {
	var r models.RouteResponse
	var routeResponse []models.RouteResponse
	for _, route := range routeList {
		r = CreateAPIResponse(route)
		routeResponse = append(routeResponse, r)
	}
	return routeResponse
}
