package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/paypermint/bridge-app-svc/models"
	"github.com/paypermint/bridge-app-svc/util"
)

//PostPaymentUpdateRequest represents struct to update payments at SAP end
type PostPaymentUpdateRequest struct {
	Records []*SapRecord `json:"Records,omitempty"`
}

// PostPaymentUpdateToSAP calls SAP api for updating payments
func (c *Client) PostPaymentUpdateToSAP(paymentUpdateRequest *PostPaymentUpdateRequest, platform string) (*SAPSuccessResponse, error) {
	jsonValue, _ := json.Marshal(paymentUpdateRequest)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/fipayconfirmationib", c.baseURL), bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Platform", platform)
	res := models.PaymentUpdateResponse{
		Records: &models.StatusRecord{},
	}
	var response *SAPSuccessResponse
	if response, err = c.sendRequestToSAP(req, res); err != nil {
		return nil, err
	}

	return response, nil
}

func getRecordParamsForJSON(r *http.Request, itemKey string) (map[string]interface{}, string, error) {
	params := make(map[string]interface{})
	var f interface{}
	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	err := dec.Decode(&f)
	if err != nil && err != io.EOF {
		return params, EmptyString, err
	}
	if f == nil {
		return params, EmptyString, nil
	}

	m, ok := f.(map[string]interface{})
	if !ok {
		return nil, EmptyString, errors.New("Can not convert to map")
	}
	value, ok := m[itemKey]
	if ok && reflect.TypeOf(value).Kind() != reflect.Slice {
		return nil, util.KeyRecords, errors.New(util.InvalidPostParameterMsg)
	}

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			params[k] = vv
		case bool:
			params[k] = strconv.FormatBool(vv)
		case json.Number:
			params[k] = vv.String()
		case map[string]interface{}:
			jsonString, err := json.Marshal(v)
			if err != nil {
				return nil, k, err
			}
			params[k] = string(jsonString)
		case []interface{}:
			switch k {
			case itemKey:
				recordItems, field, err := getRecordItemParamsForJSON(vv)
				if err != nil {
					return nil, field, err
				}
				params[k] = recordItems
			default:
				return nil, k, errors.New(util.InvalidPostParameterMsg)
			}
		default:
			return nil, k, errors.New(util.InvalidPostParameterMsg)
		}
	}
	return params, EmptyString, nil
}

func getRecordItemParamsForJSON(values []interface{}) ([]*SapRecord, string, error) {
	var records []*SapRecord

	for _, val := range values {
		// val is an interface. which contains map[string]interface{} and for note its map[string](map[string]interface{})

		switch vFurther := val.(type) {
		case map[string]interface{}:
			if field, ok := HasUnsupportedInterfaceParameters(vFurther, util.KeyCustomerNumber, util.KeyCustomerName, util.KeyCompanyCode, util.KeyItem, util.KeyAmountDue,
				util.KeyDescription, util.KeyPaymentAmount, util.KeyBankAccount, util.KeyTransactionRef); ok {
				return nil, field, errors.New(util.UnsupportedParamMsg)
			}
			record := &SapRecord{}
			customerNumber, err := GetStringInterfaceParameter(vFurther, util.KeyCustomerNumber)
			if err != nil {
				return nil, util.KeyCustomerNumber, err
			}
			delete(vFurther, util.KeyCustomerNumber)
			customerName, err := GetStringInterfaceParameter(vFurther, util.KeyCustomerName)
			if err != nil {
				return nil, util.KeyCustomerName, errors.New(util.InvalidPostParameterMsg)
			}
			delete(vFurther, util.KeyCustomerName)
			companyCode, err := GetStringInterfaceParameter(vFurther, util.KeyCompanyCode)
			if err != nil {
				return nil, util.KeyCompanyCode, err
			}
			delete(vFurther, util.KeyCompanyCode)
			description, err := GetStringInterfaceParameter(vFurther, util.KeyDescription)
			if err != nil {
				return nil, util.KeyDescription, err
			}
			delete(vFurther, util.KeyDescription)
			item, err := GetStringInterfaceParameter(vFurther, util.KeyItem)
			if err != nil {
				return nil, util.KeyItem, err
			}
			delete(vFurther, util.KeyItem)
			amountDue, err := GetStringInterfaceParameter(vFurther, util.KeyAmountDue)
			if err != nil {
				return nil, util.KeyAmountDue, err
			}
			delete(vFurther, util.KeyAmountDue)
			paymentAmount, err := GetStringInterfaceParameter(vFurther, util.KeyPaymentAmount)
			if err != nil {
				return nil, util.KeyPaymentAmount, err
			}
			delete(vFurther, util.KeyPaymentAmount)
			bankAccount, err := GetStringInterfaceParameter(vFurther, util.KeyBankAccount)
			if err != nil {
				return nil, util.KeyBankAccount, err
			}
			delete(vFurther, util.KeyBankAccount)
			transActionRef, err := GetStringInterfaceParameter(vFurther, util.KeyTransactionRef)
			if err != nil {
				return nil, util.KeyTransactionRef, err
			}
			delete(vFurther, util.KeyTransactionRef)
			record.CustomerNumber = customerNumber
			record.CustomerName = customerName
			record.CompanyCode = companyCode
			record.Description = description
			record.Item = item
			record.AmountDue = amountDue
			record.PaymentAmount = paymentAmount
			record.BankAccount = bankAccount
			record.TransactionRef = transActionRef

			jsonString, err := json.Marshal(vFurther)
			if err != nil {
				return nil, util.KeyRecords, errors.New(util.InvalidPostParameterMsg)
			}

			err = json.Unmarshal(jsonString, record)
			if err != nil {
				return nil, util.KeyRecords, errors.New(util.InvalidPostParameterMsg)
			}

			records = append(records, record)
		}
	}
	return records, "", nil
}
