package handlers

import (
	"net/http"

	"github.com/paypermint/bridge-app-svc/helpers"
	"github.com/paypermint/bridge-app-svc/util"
)

//GetAPIList renders the list of APIs provided
func GetAPIList(w http.ResponseWriter, req *http.Request) {
	routeList := GetRoutes()
	apiList := helpers.CreateAPIList(routeList)

	util.RenderJSON(appCtx, w, http.StatusOK, apiList)
}

//NotFound renders the error json if url not found
func NotFound(w http.ResponseWriter, req *http.Request) {
	util.RenderErrorJSON(appCtx, w, http.StatusNotFound, "Request URL "+req.Method+" "+req.URL.Path+" does not exist. Please check documentation", "")
}
