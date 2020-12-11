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

	params, _, _ := helpers.GetRequestParams(req, "POST")
	if field, ok := helpers.HasUnsupportedParameters(params, util.KeyFilePath); ok {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, util.UnsupportedParamMsg, field)
		return
	}

	sapClient := helpers.CreateSAPClient(req.RemoteAddr)
	ctxLogger.Info("SAP Client", "message", sapClient)
	paymentUpdateRequest := &helpers.PostPaymentUpdateRequest{
		Records: []*helpers.Record{
			&helpers.Record{
				CustomerNumber: "0001000063",
				CustomerName:   "JINDAL STEEL & POWER LIMITED",
				CompanyCode:    "MCL",
				Description:    "Proforma",
				Item:           "9101000576",
				AmountDue:      "3000.00",
				PaymentAmount:  "3000.00",
				BankAccount:    "007801022486",
				TransactionRef: "UTI0001",
			},
		},
	}
	response, err := sapClient.PostPaymentUpdateToSAP(paymentUpdateRequest)
	if err != nil {
		ctxLogger.Crit(err.Error())
		util.RenderAPIErrorJSON(appCtx, w)
		return
	}
	ctxLogger.Info("SAP Response", "message", response)

	util.RenderJSON(appCtx, w, http.StatusOK, response)
	return
}
