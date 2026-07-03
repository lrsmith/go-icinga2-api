// Package iapi provides a client for interacting with an Icinga2 Server
package iapi

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
)

// Server ... Use to be ClientConfig
type Server struct {
	Username           string
	Password           string
	BaseURL            string
	AllowUnverifiedSSL bool
	CACertFile         string
	Tries              int
	RetryDelay         time.Duration
	httpClient         *http.Client
}

func New(username, password, url string, allowUnverifiedSSL bool, caCertFile string, tries int, retryDelay time.Duration) (*Server, error) {
	return &Server{username, password, url, allowUnverifiedSSL, caCertFile, tries, retryDelay, nil}, nil
}

func (server *Server) Config(username, password, url string, allowUnverifiedSSL bool, caCertFile string, tries int, retryDelay time.Duration) (*Server, error) {
	// TODO : Add code to verify parameters
	return &Server{username, password, url, allowUnverifiedSSL, caCertFile, tries, retryDelay, nil}, nil
}

// createHttpClient defensively creates the HTTP client once
// and allow httpmock to mock the Transport attribute of the HTTP client
func (server *Server) createHttpClient() {
	if server.httpClient == nil {
		var caCertPool *x509.CertPool
		if server.CACertFile != "" {
			caCert, err := os.ReadFile(server.CACertFile)
			if err != nil {
				return
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
		}

		server.httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: server.AllowUnverifiedSSL,
					RootCAs:            caCertPool,
				},
			},
			Timeout: time.Second * 60,
		}
	}
}

func (server *Server) doRequest(method, fullURL string, body io.Reader) (*http.Response, error) {
	server.createHttpClient()

	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = io.ReadAll(body)
	}

	request, requestErr := http.NewRequest(method, fullURL, io.NopCloser(bytes.NewBuffer(bodyBytes)))
	if requestErr != nil {
		return nil, requestErr
	}

	request.SetBasicAuth(server.Username, server.Password)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	return server.httpClient.Do(request)
}

func (server *Server) Connect() error {

	response, doErr := server.doRequest("GET", server.BaseURL, nil)

	if doErr != nil || response == nil {
		server.httpClient = nil
		return doErr
	}

	defer response.Body.Close()

	return nil
}

// NewAPIRequest ...
func (server *Server) NewAPIRequest(method, APICall string, jsonString []byte) (*APIResult, error) {
	operation := func() (*APIResult, error) {
		response, doErr := server.doRequest(method, server.BaseURL+APICall, bytes.NewBuffer(jsonString))
		if doErr != nil {
			results := &APIResult{
				Code:        0,
				Status:      "Error : Request to server failed : " + doErr.Error(),
				ErrorString: doErr.Error(),
			}
			return results, backoff.Permanent(doErr)
		}
		defer response.Body.Close()

		results := &APIResult{}
		if decodeErr := json.NewDecoder(response.Body).Decode(&results); decodeErr != nil {
			return nil, backoff.Permanent(decodeErr)
		}

		if results.Code == 0 {
			results.Code = response.StatusCode
		}

		if results.Status == "" {
			results.Status = response.Status
		}

		results.ErrorString = results.Status
		if response.StatusCode == 0 {
			results.ErrorString = "Did not get a response code."
		}

		// Retry when Icinga is reloading
		// The error message returned by Icinga can have a dot ("Icinga is reloading.") or not ("Icinga is reloading")
		if results.Code == http.StatusServiceUnavailable && strings.HasPrefix(results.Status, "Icinga is reloading") {
			return results, fmt.Errorf("icinga is reloading")
		}

		return results, nil
	}

	ctx := context.Background()

	// Number of tries must be at least 1 to avoid infinite loop
	tries := uint(math.Max(float64(server.Tries), 1.0))

	return backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewConstantBackOff(server.RetryDelay)), backoff.WithMaxTries(tries))
}
