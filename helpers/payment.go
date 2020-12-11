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
	Records []*SapRecord `json:"Records,omitempty"`
}

// PostPaymentUpdateToSAP calls SAP api for updating payments
func (c *Client) PostPaymentUpdateToSAP(paymentUpdateRequest *PostPaymentUpdateRequest) (*SAPSuccessResponse, error) {
	jsonValue, _ := json.Marshal(paymentUpdateRequest)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/fipayconfirmationib", c.baseURL), bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res := models.PaymentUpdateResponse{
		Records: &models.StatusRecord{},
	}
	var response *SAPSuccessResponse
	if response, err = c.sendRequestToSAP(req, res); err != nil {
		return nil, err
	}

	return response, nil
}
