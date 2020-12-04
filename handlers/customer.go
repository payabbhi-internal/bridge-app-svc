package handlers

import (
	"fmt"
	"net/http"
)

//GET Operations

//SyncCustomers perfors syncing of customers between our end and at SAP end
func SyncCustomers(w http.ResponseWriter, req *http.Request) {
	fmt.Println("coming at this route")
}
