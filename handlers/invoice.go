package handlers

import (
	"net/http"

	"github.com/paypermint/appkit"
	"github.com/paypermint/bridge-app-svc/helpers"
	"github.com/paypermint/bridge-app-svc/util"
)

//SyncInvoices perfors syncing of invoices between our end and at SAP end
func SyncInvoices(w http.ResponseWriter, r *http.Request) {
	ctxLogger := appkit.GetContextLogger(appCtx.Logger, r)
	params, _, _ := helpers.GetRequestParams(r, "POST")
	if field, ok := helpers.HasUnsupportedParameters(params, util.KeyCount, util.KeyFrom, util.KeyTo, util.KeySkip, util.KeyBillingMethod,
		util.KeyEmail, util.KeyInvoiceNo, util.KeySubscriptionID, util.KeyDueDateFrom, util.KeyDueDateTo, util.KeyMerchantInvoiceID, util.KeyInvoiceCategory, util.KeyInvoiceIDs); ok {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, util.UnsupportedParamMsg, field)
		return
	}
	ctxLogger.Info("inside SyncInvoices")

	basicAuthCreds, bearerTokenCreds, err := helpers.GetCredentialsFromRequestHeader(r)
	if err != nil {
		ctxLogger.Crit(err.Error())
		util.RenderAPIErrorJSON(appCtx, w)
		return
	}

	c := helpers.NewClient(basicAuthCreds, bearerTokenCreds, r.RemoteAddr)
	searchInvoicesRequest := &helpers.SearchInvoicesRequest{
		MerchantInvoiceID: params[util.KeyMerchantInvoiceID],
	}

	list, err := c.GetInvoices(searchInvoicesRequest)
	if err != nil {
		ctxLogger.Crit(err.Error())
		util.RenderAPIErrorJSON(appCtx, w)
		return
	}

	util.RenderJSON(appCtx, w, http.StatusOK, list)

}
