package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	PayabbhiBaseURL = "https://0.0.0.0:50091/api/v1"
)

// Client .
type Client struct {
	accessID   string
	secretKey  string
	baseURL    string
	remoteAddr string
	HTTPClient *http.Client
}

// NewClient creates client with given API keys
func NewClient(accessID, secretKey, remoteAddr string) *Client {
	return &Client{
		accessID:  accessID,
		secretKey: secretKey,
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

// Content-type and body should be already added to req
func (c *Client) sendRequestToPayabbhi(req *http.Request, v interface{}) error {
	fmt.Println("inside sendRequestToPayabbhi")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	// req.Header.Set("Authorization", fmt.Sprintf("Basic Auth %s %s", c.accessID, c.secretKey))
	req.SetBasicAuth(c.accessID, c.secretKey)
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
