package handlers

import (
	"net/http"

	"github.com/paypermint/appkit"
	"github.com/paypermint/bridge-app-svc/helpers"
	"github.com/paypermint/bridge-app-svc/util"
)

//POST Operations

//SyncPayments perfors syncing of payments between our end and at SAP end
func SyncPayments(w http.ResponseWriter, req *http.Request) {
	ctxLogger := appkit.GetContextLogger(appCtx.Logger, req)

	ctxLogger.Info("inside SyncPayments")

	params, field, err := helpers.GetParamsWithRecordRequest(req, util.KeyRecords)
	if err != nil {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, err.Error(), field)
		return
	}
	if field, ok := helpers.HasUnsupportedInterfaceParameters(params, util.KeyRecords, field); ok {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, util.UnsupportedParamMsg, field)
		return
	}

	//Mandatory
	recordItems, err := helpers.GetRecordsParam(params, util.KeyRecords, false)
	if err != nil {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, err.Error(), util.KeyRecords)
		return
	}

	ctxLogger.Info("recordItems: ", "message", recordItems)

	traceID := appkit.TraceIDFromHTTPRequest(req)
	vClient, err := appkit.VaultConnect(appCtx, traceID)
	if err != nil {
		ctxLogger.Crit(err.Error())
		util.RenderAPIErrorJSON(appCtx, w)
		return
	}

	// sap user credentials from vault
	userid, password, fetchErr := vClient.SAPClientCreds(helpers.GetSapUserCredsPath())
	if fetchErr != nil {
		ctxLogger.Crit(fetchErr.Error())
		util.RenderAPIErrorJSON(appCtx, w)
		return
	}
	sapClient := helpers.CreateSAPClient(req.RemoteAddr, userid, password)
	ctxLogger.Info("SAP Client", "message", sapClient)
	paymentUpdateRequest := &helpers.PostPaymentUpdateRequest{
		Records: recordItems,
	}
	platform := req.Header.Get("Platform")
	response, err := sapClient.PostPaymentUpdateToSAP(paymentUpdateRequest, platform)
	if err != nil {
		ctxLogger.Crit(err.Error())
		util.RenderAPIErrorJSON(appCtx, w)
		return
	}
	ctxLogger.Info("SAP Response", "message", response)

	util.RenderJSON(appCtx, w, http.StatusOK, response)
	return
}
