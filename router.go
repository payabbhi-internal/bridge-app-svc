package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paypermint/bridge-app-svc/handlers"
	"github.com/paypermint/bridge-app-svc/models"
)

func setRouter() *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(handlers.NotFound)
	router.HandleFunc("/apilist", handlers.GetAPIList).Methods("GET")
	apiV1 := router.PathPrefix("/bridgeapp/v1").Subrouter()
	routesList := []models.Route{
		models.Route{
			Name:        "SyncCustomersAPI",
			Methods:     []string{"POST"},
			Pattern:     "/customers",
			HandlerFunc: handlers.SyncCustomers,
		},
		models.Route{
			Name:        "SyncInvoices",
			Methods:     []string{"PUT"},
			Pattern:     "/sync_invoices",
			HandlerFunc: handlers.SyncInvoices,
		},
		models.Route{
			Name:        "SyncPayments",
			Methods:     []string{"POST"},
			Pattern:     "/payments",
			HandlerFunc: handlers.SyncPayments,
		},
	}

	for _, route := range routesList {
		apiV1.
			Methods(route.Methods...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	handlers.SetRoutes(routesList)
	return router
}
