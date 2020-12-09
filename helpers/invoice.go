package helpers

import (
	"fmt"
	"net/http"

	"github.com/paypermint/bridge-app-svc/models"
)

//SearchInvoicesRequest struct
type SearchInvoicesRequest struct {
	InvoiceIDs        []string `json:"invoice_ids"`
	MerchantInvoiceID string   `json:"merchant_invoice_id"`
}

// GetInvoices calls payabbhi api for getting invoices
func (c *Client) GetInvoices(searchInvoicesRequest *SearchInvoicesRequest) (*models.List, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/invoice_ins?merchant_invoice_id=%s", c.baseURL, searchInvoicesRequest.MerchantInvoiceID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res := models.List{}
	if err := c.sendRequestToPayabbhi(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
