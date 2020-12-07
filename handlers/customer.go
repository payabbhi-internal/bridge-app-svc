package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	fmt.Println("request: ", req)
	fmt.Println("params:", params)
	fmt.Println(req.URL)
	fmt.Println(req.Host)

	fmt.Println(req.RemoteAddr)
	customersData, err := helpers.ReadCSVFile(filePath)
	if err != nil {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, err.Error(), util.KeyFilePath)
		return
	}

	client := helpers.NewClient(accessID, secretKey, req.RemoteAddr)
	fmt.Println("client:", client)
	for index, customerData := range customersData {

		if index == 0 {
			continue
		}

		var notesJSON map[string]interface{}

		// Unmarshal or Decode the JSON to the interface.
		json.Unmarshal([]byte(customerData[15]), &notesJSON)
		createCustomerRequest := &helpers.CreateCustomerRequest{
			Name:      customerData[1],
			Email:     customerData[2],
			ContactNo: customerData[3],
			BillingAddress: &helpers.Address{
				AddressLine1: customerData[4],
				AddressLine2: customerData[5],
				City:         customerData[6],
				State:        customerData[7],
				Pin:          customerData[8],
			},
			ShippingAddress: &helpers.Address{
				AddressLine1: customerData[9],
				AddressLine2: customerData[10],
				City:         customerData[11],
				State:        customerData[12],
				Pin:          customerData[13],
			},
			Gstin:              customerData[14],
			Notes:              notesJSON,
			MerchantCustomerID: customerData[17],
			HasPortalAccess:    true,
		}
		//
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
