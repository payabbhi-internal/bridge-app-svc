package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

//CreateCustomerRequest represents struct to create customer
type CreateCustomerRequest struct {
	ProfileID          string                 `json:"profileID,omitempty"`
	Name               string                 `json:"name,omitempty"`
	Email              string                 `json:"email,omitempty"`
	ContactNo          string                 `json:"contact_no,omitempty"`
	Gstin              string                 `json:"gstin,omitempty"`
	Notes              map[string]interface{} `json:"notes,omitempty"`
	BillingAddress     *Address               `json:"billing_address,omitempty"`
	ShippingAddress    *Address               `json:"shipping_address,omitempty"`
	MerchantCustomerID string                 `json:"merchant_customer_id,omitempty"`
	BankDetails        []*BankDetail          `json:"bank_details,omitempty"`
	HasPortalAccess    bool                   `json:"has_portal_access,omitempty"`
	Label              string                 `json:"label,omitempty"`
	Env                string                 `json:"env,omitempty"`
}

//Address represents Address struct
type Address struct {
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	Pin          string `json:"pin,omitempty"`
}

//BankDetail represents BankDetail struct
type BankDetail struct {
	BankName        string `json:"bank_name,omitempty"`
	Ifsc            string `json:"ifsc,omitempty"`
	AccountNo       string `json:"account_no,omitempty"`
	AccountType     string `json:"account_type,omitempty"`
	Status          int64  `json:"status,omitempty"`
	BeneficiaryName string `json:"beneficiary_name,omitempty"`
}

// CreateCustomer calls payabbhi api for creating customer
func (c *Client) CreateCustomer(createCustomerRequest *CreateCustomerRequest) error {
	jsonValue, _ := json.Marshal(createCustomerRequest)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/customers", c.baseURL), bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if err := c.sendRequestToPayabbhi(req, nil); err != nil {
		return err
	}

	return nil
}
