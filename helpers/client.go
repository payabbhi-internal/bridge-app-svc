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

const (
	PayabbhiBaseURL = "https://localhost:50091/api/v1"
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
		baseURL:    PayabbhiBaseURL,
		remoteAddr: remoteAddr,
	}
}

type payabbhiErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type payabbhiSuccessResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

//ChangeBaseURL changes the base url field set in header
func ChangeBaseURL(req *http.Request, c *Client) {
	c.baseURL = fmt.Sprintf("https://%s/api/v1", GetDynamicHost())
}

// Content-type and body should be already added to req
func (c *Client) sendRequestToPayabbhi(req *http.Request, v interface{}) error {
	fmt.Println("inside sendRequestToPayabbhi")
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
	fmt.Println("***req.Header:", req.Header)
	fmt.Println("***authorization:", req.Header.Get("Authorization"))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Try to unmarshall into errorResponse
	if res.StatusCode != http.StatusOK {
		var errRes payabbhiErrorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	fmt.Println("resp", json.NewDecoder(res.Body))

	// Unmarshall and populate v
	fullResponse := payabbhiSuccessResponse{
		Data: v,
	}
	if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
		return err
	}

	return nil
}
