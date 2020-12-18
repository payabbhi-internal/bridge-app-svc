package helpers

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/paypermint/bridge-app-svc/util"
)

const (
	methodsObject = "methods"
	//TestEnv denotes the test env
	TestEnv = "test"
	//LiveEnv denotes the live env
	LiveEnv = "live"
	//EmptyString denotes empty string
	EmptyString = ""

	urlEncodedHeader = "application/x-www-form-urlencoded"
	jsonHeader       = "application/json"
)

const (
	reqAccessID           = "access_id"
	authError             = "Authorization Header not present"
	headerValueMissing    = "AccessID/SecretKey missing in the authorization header"
	invalidInvoiceRequest = "Error processing invoice request"
	tokenSigningAlgo      = "RS256"
	validAuthPrefix       = "Bearer "
	skanAppDomain         = "skan_app"
	b2bPortalDomain       = "b2b_portal"
)

//SapRecord represents SapRecord struct
type SapRecord struct {
	CustomerNumber string `json:"customer_number,omitempty"`
	CustomerName   string `json:"customer_name,omitempty"`
	CompanyCode    string `json:"company_code,omitempty"`
	Description    string `json:"description,omitempty"`
	Item           string `json:"item,omitempty"`
	AmountDue      string `json:"amount_due,omitempty"`
	PaymentAmount  string `json:"payment_amount,omitempty"`
	BankAccount    string `json:"bank_account,omitempty"`
	TransactionRef string `json:"transaction_ref,omitempty"`
	CustomerID     string `json:"Customer_ID,omitempty"`
}

func isJWTAuthenticationRequest(request *http.Request) (string, bool) {
	authToken := request.Header.Get("Authorization")
	return authToken, strings.HasPrefix(authToken, validAuthPrefix)
}

//GetCredentialsFromRequestHeader .
func GetCredentialsFromRequestHeader(request *http.Request) (*BasicAuthCreds, *BearerAuthCreds, error) {
	authToken, authTokenOk := isJWTAuthenticationRequest(request)
	accessID, secretKey, basicAuthOk := request.BasicAuth()

	if !basicAuthOk && !authTokenOk {
		return nil, nil, errors.New(headerValueMissing)
	}
	if !basicAuthOk {
		return nil, &BearerAuthCreds{
			token: authToken,
		}, nil
	}
	if !authTokenOk {
		return &BasicAuthCreds{
			accessID:  accessID,
			secretKey: secretKey,
		}, nil, nil
	}
	return &BasicAuthCreds{
			accessID:  accessID,
			secretKey: secretKey,
		}, &BearerAuthCreds{
			token: authToken,
		}, nil
}

//GetRequestParams takes the request params and depending on the method forms the params
func GetRequestParams(r *http.Request, method string) (map[string]string, string, error) {
	switch method {
	case "POST", "PUT":
		return getPostParams(r)
	case "GET":
		params := map[string]string{}
		for key := range r.URL.Query() {
			params[key] = r.URL.Query().Get(key)
		}
		return params, EmptyString, nil
	default:
		return nil, "method", errors.New("Unknown request method")
	}
}

func getPostParams(r *http.Request) (map[string]string, string, error) {
	params := make(map[string]string)
	var field string
	errPayload := errors.New("Error in payload formation")
	contentType := getContentTypeFromRequest(r.Header)
	switch contentType {
	case jsonHeader:
		params, field, errPayload = getPostParamsForJSON(r)
	// case urlEncodedHeader:
	// 	params, field, errPayload = getPostParamsForURLEncoded(r)
	default:
		return nil, "Content-Type", errPayload
	}
	if errPayload != nil {
		return nil, field, errPayload
	}
	return params, field, nil
}

func getContentTypeFromRequest(header http.Header) string {
	if header["Content-Type"] != nil {
		return header["Content-Type"][0]
	}
	return jsonHeader
}

func getPostParamsForJSON(r *http.Request) (map[string]string, string, error) {
	params := map[string]string{}
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

	value, ok := m[util.KeyNotes]
	if ok && reflect.TypeOf(value).Kind() != reflect.Map {
		return nil, util.KeyNotes, errors.New(util.InvalidPostParameterMsg)
	}

	value, ok = m[util.KeyShippingAddress]
	if ok && reflect.TypeOf(value).Kind() != reflect.Map {
		return nil, util.KeyShippingAddress, errors.New(util.InvalidPostParameterMsg)
	}

	value, ok = m[util.KeyBillingAddress]
	if ok && reflect.TypeOf(value).Kind() != reflect.Map {
		return nil, util.KeyBillingAddress, errors.New(util.InvalidPostParameterMsg)
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
			jsonString, err := json.Marshal(v)
			if err != nil {
				return nil, k, err
			}
			params[k] = string(jsonString)
		default:
			return nil, k, fmt.Errorf("Unknown type for %s", vv)
		}
	}
	return params, EmptyString, nil
}

//HasUnsupportedParameters returns true if the data map contains a key not present in the list of keys
func HasUnsupportedParameters(data map[string]string, keys ...string) (string, bool) {
	expected := map[string]bool{}
	for _, key := range keys {
		expected[key] = true
	}
	for key := range data {
		if !expected[key] {
			return key, true
		}
	}
	return EmptyString, false
}

//ReadCSVFromURL returns data from file location at a given url
func ReadCSVFromURL(url string) ([][]string, error) {
	resp, err := http.Get("http://" + url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ';'
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

//ReadCSVFile returns data from file location
func ReadCSVFile(filePath string) ([][]string, error) {
	data := make([][]string, 0)
	csvFile, err := os.Open(filePath)
	if err != nil {
		return data, err
	}
	r := csv.NewReader(csvFile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println(record)
		data = append(data, record)
	}
	return data, nil
}

//GetStringParam returns the value of key in params
func GetStringParam(params map[string]string, key string) (string, error) {
	if value, ok := params[key]; ok {
		if value == EmptyString {
			return EmptyString, errors.New(util.InvalidPostParameterMsg)
		}
		return value, nil
	}
	return EmptyString, errors.New(util.MissingMandatoryField)
}

//GetOptionalStringParam returns the value of key in params
func GetOptionalStringParam(params map[string]string, key string) (string, error) {
	if value, ok := params[key]; ok {
		if value == EmptyString {
			return EmptyString, errors.New(util.InvalidPostParameterMsg)
		}
		return value, nil
	}
	return EmptyString, nil
}

//GetStringInterfaceParam returns the value of key in params
func GetStringInterfaceParam(params map[string]interface{}, key string, optional bool) (string, error) {
	if value, ok := params[key]; ok {
		if value == EmptyString && !optional {
			return EmptyString, errors.New(util.InvalidPostParameterMsg)
		}

		switch val := value.(type) {
		case json.Number:
			return val.String(), nil
		case string:
			return val, nil
		case int:
			return strconv.Itoa(val), nil
		case int64:
			return strconv.FormatInt(val, 10), nil
		case float64:
			return strconv.FormatFloat(val, 'f', 0, 64), nil
		default:
			if !optional {
				return EmptyString, errors.New(util.InvalidPostParameterMsg)
			}
		}
	}
	if optional {
		return EmptyString, nil
	}

	return EmptyString, errors.New(util.MissingMandatoryField)
}

//GetAmountParamInPaisa returns the amount field's value in paisa
func GetAmountParamInPaisa(params map[string]interface{}, key string, optional, isPositive bool) (int64, error) {
	if value, ok := params[key]; ok {
		switch val := value.(type) {
		case json.Number:
			valJSON := val.String()
			intValue, err := strconv.ParseInt(valJSON, 10, 64)
			if err != nil {
				return 0, errors.New(util.InvalidPostParameterMsg)
			}
			if isPositive && intValue <= 0 {
				return 0, errors.New(util.InvalidPostParameterMsg)
			}
			return intValue * 100, nil
		case string:
			floatValue, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
			if err != nil {
				return 0, errors.New(util.InvalidPostParameterMsg)
			}
			intValue := int64(floatValue * 100)
			if isPositive && intValue < 0 {
				return 0, errors.New(util.InvalidPostParameterMsg)
			}
			return intValue, nil
		case int64:
			if isPositive && val <= 0 {
				return 0, errors.New(util.InvalidPostParameterMsg)
			}
			return val * 100, nil
		case float64:
			intValue := int64(val * 100)
			if isPositive && intValue <= 0 {
				return 0, errors.New(util.InvalidPostParameterMsg)
			}
			return intValue, nil
		default:
			return 0, errors.New(util.InvalidPostParameterMsg)
		}
	}
	if optional {
		return 0, nil
	}
	return 0, errors.New(util.MissingMandatoryField)
}

//HasUnsupportedInterfaceParameters returns true if the data map contains a key not present in the list of keys
func HasUnsupportedInterfaceParameters(data map[string]interface{}, keys ...string) (string, bool) {
	expected := map[string]bool{}
	for _, key := range keys {
		expected[key] = true
	}
	for key := range data {
		if !expected[key] {
			return key, true
		}
	}
	return EmptyString, false
}

//GetStringInterfaceParameter returns the value of key in params
//TODO refactor
func GetStringInterfaceParameter(params map[string]interface{}, key string) (string, error) {
	if value, ok := params[key]; ok {
		if value == EmptyString || (key == util.KeyCurrency && params[key] != "INR") {
			return EmptyString, errors.New(util.InvalidPostParameterMsg)
		}
		return value.(string), nil
	}
	return EmptyString, errors.New(util.MissingMandatoryField)
}

//GetParamsWithRecordRequest returns the api request params as an map of key to interface
func GetParamsWithRecordRequest(r *http.Request, itemKey string) (map[string]interface{}, string, error) {
	params := make(map[string]interface{})
	var field string
	var err error
	errPayload := errors.New("Error in payload formation")
	contentType := getContentTypeFromRequest(r.Header)

	switch contentType {
	case jsonHeader:
		if itemKey == util.KeyRecords {
			params, field, err = getRecordParamsForJSON(r, itemKey)
		}

	default:
		return nil, field, errPayload
	}
	if err != nil {
		return nil, field, err
	}
	return params, field, nil
}

//GetRecordsParam creates a record request from param map
func GetRecordsParam(params map[string]interface{}, recordKey string, optional bool) ([]*SapRecord, error) {
	var recordRequests []*SapRecord
	if recordItems, ok := params[recordKey]; ok {
		switch vv := recordItems.(type) {
		case []*SapRecord:
			if len(vv) == 0 {
				return nil, errors.New(util.InvalidPostParameterMsg)
			}
			for _, record := range vv {
				recordRequests = append(recordRequests, &SapRecord{
					CustomerNumber: record.CustomerNumber,
					CustomerName:   record.CustomerName,
					CompanyCode:    record.CompanyCode,
					Description:    record.Description,
					Item:           record.Item,
					AmountDue:      record.AmountDue,
					PaymentAmount:  record.PaymentAmount,
					BankAccount:    record.BankAccount,
					TransactionRef: record.TransactionRef,
				})
			}
		default:
			return nil, fmt.Errorf("Unknown type for %s", vv)
		}
	} else {
		if !optional {
			return nil, errors.New(util.MissingMandatoryField)
		}
	}
	return recordRequests, nil
}
