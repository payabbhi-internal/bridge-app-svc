package helpers

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client .
type Client struct {
	basicAuthCreds *BasicAuthCreds
	tokenAuthCreds *BearerAuthCreds
	baseURL        string
	remoteAddr     string
	HTTPClient     *http.Client
}

//BasicAuthCreds .
type BasicAuthCreds struct {
	accessID  string
	secretKey string
}

//BearerAuthCreds .
type BearerAuthCreds struct {
	token       string
	environment string
}

// NewClient creates client with given API keys
func NewClient(basicAuthCreds *BasicAuthCreds, bearerTokenCreds *BearerAuthCreds, remoteAddr string) *Client {
	return &Client{
		basicAuthCreds: basicAuthCreds,
		tokenAuthCreds: bearerTokenCreds,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		baseURL:    fmt.Sprintf("https://%s/api/v1", GetDynamicHost()),
		remoteAddr: remoteAddr,
	}
}

// CreateSAPClient creates new SAP Client with given username and password
func CreateSAPClient(remoteAddr, userid, password string) *Client {
	return &Client{
		basicAuthCreds: &BasicAuthCreds{
			accessID:  userid,
			secretKey: password,
		},
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		baseURL:    fmt.Sprintf("http://%s", GetSapURL()),
		remoteAddr: remoteAddr,
	}
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SAPSuccessResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// Content-type and body should be already added to req
func (c *Client) sendRequestToPayabbhi(req *http.Request, v interface{}) error {
	if strings.Contains(req.Host, "localhost") {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	if c.basicAuthCreds != nil {
		req.SetBasicAuth(c.basicAuthCreds.accessID, c.basicAuthCreds.secretKey)
	}
	if c.tokenAuthCreds != nil {
		req.Header.Add("Authorization", c.tokenAuthCreds.token)
		// req.Header.Add("env", c.tokenAuthCreds.environment)
	}
	req.RemoteAddr = ""

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Try to unmarshall into errorResponse
	if res.StatusCode != http.StatusOK {
		var errRes ErrorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	// Unmarshall and populate v
	fullResponse := SAPSuccessResponse{
		Data: v,
	}
	if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
		return err
	}

	return nil
}

// Content-type and body should be already added to req
func (c *Client) sendRequestToSAP(req *http.Request, v interface{}) (*SAPSuccessResponse, error) {
	req.Header.Set("Accept", "application/json; charset=utf-8")
	if c.basicAuthCreds != nil {
		req.SetBasicAuth(c.basicAuthCreds.accessID, c.basicAuthCreds.secretKey)
	}
	req.RemoteAddr = ""

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		errRes := ErrorResponse{Code: res.StatusCode}
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return nil, errors.New(errRes.Message)
		}

		return nil, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	// Unmarshall and populate v
	fullResponse := SAPSuccessResponse{
		Data: v,
		Code: res.StatusCode,
	}
	if err = json.NewDecoder(res.Body).Decode(&fullResponse.Data); err != nil {
		return nil, err
	}

	return &fullResponse, nil
}
