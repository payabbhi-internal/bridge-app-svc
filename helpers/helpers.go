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
	resp, err := http.Get("file:/" + url)
	if err != nil {
		fmt.Println("Err::  ", err)
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ';'
	data, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Err1::  ", err)
		return nil, err
	}

	return data, nil
}

//ReadCSVFile returns data from file location
func ReadCSVFile(filePath string) ([][]string, error) {
	data := make([][]string, 0)
	csvFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error coming here")
		return data, err
	}
	r := csv.NewReader(csvFile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			fmt.Println("Ending now")
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
