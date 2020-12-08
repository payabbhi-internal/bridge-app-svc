package handlers

import (
	"fmt"
	"net/http"

	"github.com/paypermint/appkit"
	"github.com/paypermint/bridge-app-svc/helpers"
	"github.com/paypermint/bridge-app-svc/util"
)

func SyncInvoices(w http.ResponseWriter, r *http.Request) {
	fmt.Println("**inside SyncInvoices")
	ctxLogger := appkit.GetContextLogger(appCtx.Logger, r)
	traceID := appkit.TraceIDFromHTTPRequest(r)
	params, _, _ := helpers.GetRequestParams(r, "POST")
	if field, ok := helpers.HasUnsupportedParameters(params, util.KeyCount, util.KeyFrom, util.KeyTo, util.KeySkip, util.KeyBillingMethod,
		util.KeyEmail, util.KeyInvoiceNo, util.KeySubscriptionID, util.KeyDueDateFrom, util.KeyDueDateTo, util.KeyMerchantInvoiceID, util.KeyInvoiceCategory, util.KeyInvoiceIDs); ok {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, util.UnsupportedParamMsg, field)
		return
	}
	fmt.Println(ctxLogger)
	fmt.Println(traceID)

	fmt.Println("***who:", r.Header.Get("who"))

	basicAuthCreds, brererTokenCreds, err := helpers.GetCredentialsFromRequestHeader(r)
	if err != nil {
		ctxLogger.Crit(err.Error())
		util.RenderAPIErrorJSON(appCtx, w)
		return
	}

	c := helpers.NewClient(basicAuthCreds, brererTokenCreds, r.RemoteAddr)

	searchInvoicesRequest := &helpers.SearchInvoicesRequest{
		MerchantInvoiceID: params[util.KeyMerchantInvoiceID],
	}

	ip := appkit.GetIP(r)
	fmt.Println("**ip:", ip)
	list, err := c.GetInvoices(searchInvoicesRequest)
	if err != nil {
		ctxLogger.Crit(err.Error())
		util.RenderAPIErrorJSON(appCtx, w)
		return
	}

	util.RenderJSON(appCtx, w, http.StatusOK, list)

}
