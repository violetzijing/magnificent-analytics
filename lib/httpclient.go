package lib

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// HTTPClientInterface packages function as an interface
type HTTPClientInterface interface {
	DispatchHTTPRequest(url string, method string, body io.Reader, timeout time.Duration, tryTimes int) ([]byte, int, error)
}

type httpClient struct{}

var GlobalHTTPClient HTTPClientInterface = &httpClient{}

// constant variables for sending HTTP request
const (
	Timeout          = 3000 * time.Millisecond
	TryTimes         = 2
	TimeSleep        = 100
	DefaultErrorCode = 0

	HTTPGet = "Get"
)

// CreateRequest will generate url according to passed parameters
func CreateRequest(url string, method string, body io.Reader) (*http.Request, error) {
	var (
		req *http.Request
		err error
	)
	switch method {
	case HTTPGet:
		req, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid HTTP method")
	}

	return req, nil
}

// DispatchHTTPRequest
func (c *httpClient) DispatchHTTPRequest(url string, method string, body io.Reader, timeout time.Duration, tryTimes int) ([]byte, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	netClient := &http.Client{Timeout: timeout}

	var respBody []byte
	req, err := CreateRequest(url, method, body)
	if err != nil {
		return nil, DefaultErrorCode, err
	}
	resp, err := netClient.Do(req)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp == nil {
		return nil, DefaultErrorCode, errors.New("Failed to fetch result")
	}
	defer func() {
		if resp != nil {
			if err := resp.Body.Close(); err != nil {
				log.Errorf("Close Connect Error: %s", err.Error())
			}
		}
	}()

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	if resp.StatusCode/100 != 2 { //200 ~ 299 are success code
		return nil, resp.StatusCode, errors.New(fmt.Sprintf("Get request error: %s, %s", resp.Status, string(respBody)))
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, DefaultErrorCode, errors.New("timeout")
	}

	return respBody, resp.StatusCode, nil
}
