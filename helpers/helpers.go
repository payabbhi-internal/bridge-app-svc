package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"

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

func GetCredentialsFromRequestHeader(request *http.Request) (string, string, error) {
	accessID, secretKey, ok := request.BasicAuth()
	if !ok {
		return "", "", errors.New(authError)
	}
	if accessID == "" || secretKey == "" {
		return "", "", errors.New(headerValueMissing)
	}
	return accessID, secretKey, nil
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
