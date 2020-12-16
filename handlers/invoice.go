package handlers

import (
	"net/http"

	"github.com/paypermint/bridge-app-svc/helpers"
	"github.com/paypermint/bridge-app-svc/util"
)

//SyncInvoices performs syncing of invoices between payabbhi & other system
func SyncInvoices(w http.ResponseWriter, req *http.Request) {
	syncWith := req.Header.Get("sync_with")

	if syncWith == util.SyncWithSAP {
		helpers.SyncInvoicesWithSAP(w, req, appCtx)
	}

}
