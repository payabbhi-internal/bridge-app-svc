package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/paypermint/bridge-app-svc/models"
)

//PostPaymentUpdateRequest represents struct to update payments at SAP end
type PostPaymentUpdateRequest struct {
	Records []*Record `json:"Records,omitempty"`
}

//Record represents Record struct
type Record struct {
	CustomerNumber string `json:"customer_number,omitempty"`
	CustomerName   string `json:"customer_name,omitempty"`
	CompanyCode    string `json:"company_code,omitempty"`
	Description    string `json:"description,omitempty"`
	Item           string `json:"item,omitempty"`
	AmountDue      string `json:"amount_due,omitempty"`
	PaymentAmount  string `json:"payment_amount,omitempty"`
	BankAccount    string `json:"bank_account,omitempty"`
	TransactionRef string `json:"transaction_ref,omitempty"`
}

// PostPaymentUpdateToSAP calls SAP api for updating payments
func (c *Client) PostPaymentUpdateToSAP(paymentUpdateRequest *PostPaymentUpdateRequest) (*PayabbhiSuccessResponse, error) {
	jsonValue, _ := json.Marshal(paymentUpdateRequest)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/fipayconfirmationib", c.baseURL), bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res := models.PaymentUpdateResponse{
		Records: &models.StatusRecord{},
	}
	var response *PayabbhiSuccessResponse
	if response, err = c.sendRequestToSAP(req, res); err != nil {
		return nil, err
	}

	return response, nil
}
