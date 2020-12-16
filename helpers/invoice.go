package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/paypermint/appkit"
	"github.com/paypermint/bridge-app-svc/util"
)

//GetInvoicesFromSapRequest represents struct to get invoices from SAP system
type GetInvoicesFromSapRequest struct {
	Records []*SapRecord `json:"Records,omitempty"`
}

//GetInvoicesFromSapResponse represents struct of invoices received from SAP system
type GetInvoicesFromSapResponse struct {
	Records []*SapRecord `json:"Records,omitempty"`
}

//CreateOrUpdatePayabbhiInvoiceRequest represents struct to create or update invoice Request
type CreateOrUpdatePayabbhiInvoiceRequest struct {
	CustomerID             string                 `json:"customer_id,omitempty"`
	MerchantInvoiceID      string                 `json:"merchant_invoice_id,omitempty"`
	BillingMethod          string                 `json:"billing_method,omitempty"`
	PlaceOfSupply          string                 `json:"place_of_supply,omitempty"`
	Description            string                 `json:"description,omitempty"`
	Currency               string                 `json:"currency,omitempty"`
	NotifyBy               string                 `json:"notify_by,omitempty"`
	InvoiceNo              string                 `json:"invoice_no,omitempty"`
	TermsConditions        string                 `json:"terms_conditions,omitempty"`
	CustomerNotes          string                 `json:"customer_notes,omitempty"`
	CustomerNotificationBy string                 `json:"customer_notification_by,omitempty"`
	Notes                  map[string]interface{} `json:"notes,omitempty"`
	LineItems              []*LineItem            `json:"line_items,omitempty"`
	AmountDue              int64                  `json:"amount_due,omitempty"`
	PartialPaymentMode     bool                   `json:"partial_payment_mode,omitempty"`
	Label                  string                 `json:"label,omitempty"`
}

//LineItem represents LineItem struct
type LineItem struct {
	Id                    string `json:"id,omitempty"`
	Name                  string `json:"name,omitempty"`
	Amount                int64  `json:"amount,omitempty"`
	Currency              string `json:"currency,omitempty"`
	MerchantInvoiceItemId string `json:"merchant_invoice_item_id,omitempty"`
}

// GetInvoicesFromSap calls SAP api for fetching invoices
func (c *Client) GetInvoicesFromSap(getInvoicesFromSapRequest *GetInvoicesFromSapRequest) (*SAPSuccessResponse, error) {
	jsonValue, _ := json.Marshal(getInvoicesFromSapRequest)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/fipaycollectionib", c.baseURL), bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res := GetInvoicesFromSapResponse{
		Records: []*SapRecord{},
	}
	var response *SAPSuccessResponse
	if response, err = c.sendRequestToSAP(req, res); err != nil {
		return nil, err
	}

	return response, nil
}

// CreateOrUpdatePayabbhiInvoice calls payabbhi api for creating or updating invoice
func (c *Client) CreateOrUpdatePayabbhiInvoice(createOrUpdatePayabbhiInvoiceRequest *CreateOrUpdatePayabbhiInvoiceRequest) error {
	jsonValue, _ := json.Marshal(createOrUpdatePayabbhiInvoiceRequest)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/invoice_ins", c.baseURL), bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if err := c.sendRequestToPayabbhi(req, nil); err != nil {
		return err
	}

	return nil
}

func toGetInvoicesFromSapRequest(merchantCustomerID string) *GetInvoicesFromSapRequest {
	return &GetInvoicesFromSapRequest{
		Records: []*SapRecord{
			{
				CustomerID: merchantCustomerID,
			},
		},
	}
}

func toCreateOrUpdatePayabbhiInvoiceRequest(w http.ResponseWriter, appCtx *appkit.AppContext, params map[string]string, sapInvoiceObj map[string]interface{}) (*CreateOrUpdatePayabbhiInvoiceRequest, error) {
	//Mandatory
	customerID, err := GetStringParam(params, util.KeyCustomerID)
	if err != nil {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, err.Error(), util.KeyCustomerID)
		return nil, err
	}

	//optional
	description, err := GetStringInterfaceParam(sapInvoiceObj, util.KeySapDescription, true)
	if err != nil {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, err.Error(), util.KeySapDescription)
		return nil, err
	}

	//optional
	item, err := GetStringInterfaceParam(sapInvoiceObj, util.KeySapItem, true)
	if err != nil {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, err.Error(), util.KeySapItem)
		return nil, err
	}

	//optional
	label, err := GetStringInterfaceParam(sapInvoiceObj, util.KeySapCompanyCode, true)
	if err != nil {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, err.Error(), util.KeySapCompanyCode)
		return nil, err
	}

	//Mandatory
	amountDue, err := GetAmountParamInPaisa(sapInvoiceObj, util.KeySapAmountDue, false, true)
	if err != nil {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, err.Error(), util.KeySapAmountDue)
	}

	return &CreateOrUpdatePayabbhiInvoiceRequest{
		CustomerID:         customerID,
		MerchantInvoiceID:  item,
		Description:        description,
		AmountDue:          amountDue,
		PartialPaymentMode: true,
		Currency:           util.CurrencyINR,
		Label:              label,
		LineItems: []*LineItem{
			{
				MerchantInvoiceItemId: item,
				Name:                  fmt.Sprintf("%s_item", description),
				Currency:              util.CurrencyINR,
				Amount:                amountDue,
			},
		},
	}, nil
}

// SyncInvoicesWithSAP performs syncing of invoices between payabbhi & SAP system
func SyncInvoicesWithSAP(w http.ResponseWriter, req *http.Request, appCtx *appkit.AppContext) {
	ctxLogger := appkit.GetContextLogger(appCtx.Logger, req)
	traceID := appkit.TraceIDFromHTTPRequest(req)
	vClient, err := appkit.VaultConnect(appCtx, traceID)
	if err != nil {
		ctxLogger.Crit(err.Error())
		util.RenderAPIErrorJSON(appCtx, w)
		return
	}

	// sap user credentials from vault
	userid, password, fetchErr := vClient.SAPClientCreds(GetSapUserCredsPath())
	if fetchErr != nil {
		ctxLogger.Crit(fetchErr.Error())
		util.RenderAPIErrorJSON(appCtx, w)
		return
	}
	sapClient := CreateSAPClient(req.RemoteAddr, userid, password)
	ctxLogger.Info("SAP Client", "message", sapClient)

	params, _, _ := GetRequestParams(req, "PUT")
	if field, ok := HasUnsupportedParameters(params, util.KeyMerchantCustomerID, util.KeyCustomerID); ok {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, util.UnsupportedParamMsg, field)
		return
	}

	//Mandatory
	merchantCustomerID, err := GetStringParam(params, util.KeyMerchantCustomerID)
	if err != nil {
		util.RenderErrorJSON(appCtx, w, http.StatusBadRequest, err.Error(), util.KeyMerchantCustomerID)
		return
	}

	getInvoicesFromSapRequest := toGetInvoicesFromSapRequest(merchantCustomerID)
	ctxLogger.Info("calling SAP api for fetching invoices", "request", getInvoicesFromSapRequest)
	sapResponse, err := sapClient.GetInvoicesFromSap(getInvoicesFromSapRequest)
	if err != nil {
		ctxLogger.Crit(err.Error())
		util.RenderAPIErrorJSON(appCtx, w)
		return
	}
	ctxLogger.Info("SAP Invoice Response", "response", sapResponse)

	pushInvoicesToPayabbhi(w, req, params, appCtx, sapResponse)

	util.RenderJSON(appCtx, w, http.StatusOK, sapResponse)
	return
}

func pushInvoicesToPayabbhi(w http.ResponseWriter, req *http.Request, params map[string]string, appCtx *appkit.AppContext, sapResponse *SAPSuccessResponse) {
	ctxLogger := appkit.GetContextLogger(appCtx.Logger, req)
	basicAuthCreds, bearerTokenCreds, err := GetCredentialsFromRequestHeader(req)
	if err != nil {
		ctxLogger.Crit(err.Error())
		util.RenderAPIErrorJSON(appCtx, w)
		return
	}

	sapData, ok := sapResponse.Data.(map[string]interface{})
	if !ok {
		ctxLogger.Crit("invalid format received,Unable to parse")
		util.RenderAPIErrorJSON(appCtx, w)
	}

	if sapInvoicesData, present := sapData["Records"]; present {
		if sapInvoices, ok := sapInvoicesData.([]interface{}); ok {
			payabbhiClient := NewClient(basicAuthCreds, bearerTokenCreds, req.RemoteAddr)
			for _, sapInvoice := range sapInvoices {
				if sapInvoiceObj, ok := sapInvoice.(map[string]interface{}); ok {
					createOrUpdatePayabbhiInvoiceRequest, err := toCreateOrUpdatePayabbhiInvoiceRequest(w, appCtx, params, sapInvoiceObj)
					if err != nil {
						ctxLogger.Crit(err.Error())
						util.RenderAPIErrorJSON(appCtx, w)
						return
					}
					ctxLogger.Info("calling payabbhi CreateOrUpdateInvoice api", "request", createOrUpdatePayabbhiInvoiceRequest)
					err = payabbhiClient.CreateOrUpdatePayabbhiInvoice(createOrUpdatePayabbhiInvoiceRequest)
					if err != nil {
						ctxLogger.Crit(err.Error())
						util.RenderAPIErrorJSON(appCtx, w)
						return
					}
					ctxLogger.Info("Payabbhi CreateOrUpdateInvoice Response", "response", sapResponse)
				}
			}
		}
	}
}
