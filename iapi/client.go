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

func buildHttpClient(allowUnverifiedSSL bool, caCertFile string) (*http.Client, error) {
	var caCertPool *x509.CertPool
	if caCertFile != "" {
		caCert, err := os.ReadFile(caCertFile)
		if err != nil {
			return nil, err
		}
		caCertPool = x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: allowUnverifiedSSL,
				RootCAs:            caCertPool,
			},
		},
		Timeout: time.Second * 60,
	}, nil
}

func New(username, password, url string, allowUnverifiedSSL bool, caCertFile string, tries int, retryDelay time.Duration) (*Server, error) {
	client, err := buildHttpClient(allowUnverifiedSSL, caCertFile)
	if err != nil {
		return nil, err
	}
	return &Server{username, password, url, allowUnverifiedSSL, caCertFile, tries, retryDelay, client}, nil
}

func (server *Server) Config(username, password, url string, allowUnverifiedSSL bool, caCertFile string, tries int, retryDelay time.Duration) (*Server, error) {
	client, err := buildHttpClient(allowUnverifiedSSL, caCertFile)
	if err != nil {
		return nil, err
	}
	return &Server{username, password, url, allowUnverifiedSSL, caCertFile, tries, retryDelay, client}, nil
}

func (server *Server) doRequest(ctx context.Context, method, fullURL string, body io.Reader) (*http.Response, error) {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = io.ReadAll(body)
	}

	request, requestErr := http.NewRequestWithContext(ctx, method, fullURL, io.NopCloser(bytes.NewBuffer(bodyBytes)))
	if requestErr != nil {
		return nil, requestErr
	}

	request.SetBasicAuth(server.Username, server.Password)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	return server.httpClient.Do(request)
}

func (server *Server) Connect(ctx context.Context) error {
	response, doErr := server.doRequest(ctx, "GET", server.BaseURL, nil)
	if doErr != nil || response == nil {
		return doErr
	}
	defer response.Body.Close()

	return nil
}

// NewAPIRequest ...
func (server *Server) NewAPIRequest(ctx context.Context, method, APICall string, jsonString []byte, dest any) (*APIResult, error) {
	operation := func() (*APIResult, error) {
		response, doErr := server.doRequest(ctx, method, server.BaseURL+APICall, bytes.NewBuffer(jsonString))
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

		// Decode Results into destination pointer if provided and Results is not empty
		if dest != nil && len(results.Results) > 0 {
			if unmarshalErr := json.Unmarshal(results.Results, dest); unmarshalErr != nil {
				return nil, backoff.Permanent(unmarshalErr)
			}
		}

		// Retry when Icinga is reloading
		// The error message returned by Icinga can have a dot ("Icinga is reloading.") or not ("Icinga is reloading")
		if results.Code == http.StatusServiceUnavailable && strings.HasPrefix(results.Status, "Icinga is reloading") {
			return results, fmt.Errorf("icinga is reloading")
		}

		return results, nil
	}

	// Number of tries must be at least 1 to avoid infinite loop
	tries := uint(math.Max(float64(server.Tries), 1.0))

	return backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewConstantBackOff(server.RetryDelay)), backoff.WithMaxTries(tries))
}
