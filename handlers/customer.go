package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paypermint/appkit"
	"github.com/paypermint/bridge-app-svc/helpers"
	"github.com/paypermint/bridge-app-svc/util"
)

//GET Operations

//SyncCustomers perfors syncing of customers between our end and at SAP end
func SyncCustomers(w http.ResponseWriter, req *http.Request) {
	ctxLogger := appkit.GetContextLogger(appCtx.Logger, req)

	fmt.Println("***who:", req.Header.Get("who"))

	accessID, secretKey, err := helpers.GetCredentialsFromRequestHeader(req)
	if err != nil {
		ctxLogger.Crit(err.Error())
		util.RenderAPIErrorJSON(appCtx, w)
		return
	}
	params, _, _ := helpers.GetRequestParams(req, "POST")
	if field, ok := helpers.HasUnsupportedParameters(params, util.KeyFilePath); ok {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, util.UnsupportedParamMsg, field)
		return
	}

	//Mandatory
	filePath, err := helpers.GetStringParam(params, util.KeyFilePath)
	if err != nil {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, err.Error(), util.KeyFilePath)
		return
	}

	customersData, err := helpers.ReadCSVFile(filePath)
	if err != nil {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, err.Error(), util.KeyFilePath)
		return
	}

	customersDataFromS3, err := helpers.GetS3File(ctxLogger, "", "", "", "customers.csv", mux.Vars(req)["id"])
	if err != nil {
		ctxLogger.Error(err.Error())
		return
	}

	for _, customerData := range customersDataFromS3 {
		fmt.Println(customerData)
	}

	client := helpers.NewClient(accessID, secretKey, req.RemoteAddr)
	for index, customerData := range customersData {

		if index == 0 {
			continue
		}

		var notesJSON map[string]interface{}

		// Unmarshal or Decode the JSON to the interface.
		json.Unmarshal([]byte(customerData[14]), &notesJSON)
		createCustomerRequest := &helpers.CreateCustomerRequest{
			Name:      customerData[0],
			Email:     customerData[1],
			ContactNo: customerData[2],
			BillingAddress: &helpers.Address{
				AddressLine1: customerData[3],
				AddressLine2: customerData[4],
				City:         customerData[5],
				State:        customerData[6],
				Pin:          customerData[7],
			},
			ShippingAddress: &helpers.Address{
				AddressLine1: customerData[8],
				AddressLine2: customerData[9],
				City:         customerData[10],
				State:        customerData[11],
				Pin:          customerData[12],
			},
			Gstin:              customerData[13],
			Notes:              notesJSON,
			MerchantCustomerID: customerData[16],
			HasPortalAccess:    true,
		}
		err := client.CreateCustomer(createCustomerRequest)
		if err != nil {
			ctxLogger.Crit(err.Error())
			util.RenderAPIErrorJSON(appCtx, w)
			return
		}
		fmt.Println(customerData)
	}
	util.RenderJSON(appCtx, w, http.StatusOK, "OK!")
	return
}
