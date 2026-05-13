package iapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"

	"github.com/cenkalti/backoff/v5"
)

// GetHost ...
func (server *Server) GetHost(hostname string) ([]HostStruct, error) {

	var hosts []HostStruct

	results, err := server.NewAPIRequest(http.MethodGet, "/objects/hosts/"+hostname, nil)
	if err != nil {
		return nil, err
	}

	// Contents of the results is an interface object. Need to convert it to json first.
	jsonStr, marshalErr := json.Marshal(results.Results)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// then the JSON can be pushed into the appropriate struct.
	// Note : Results is a slice so much push into a slice.

	if unmarshalErr := json.Unmarshal(jsonStr, &hosts); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return hosts, err

}

// CreateHost creates a host.
// When a context deadline is exceeded, wait for the host to be created if a number of tries is defined.
func (server *Server) CreateHost(hostname, address, address6 string, checkCommand string, variables map[string]interface{}, templates []string, groups []string, zone string) (hosts []HostStruct, err error) {

	var newAttrs HostAttrs
	newAttrs.Address = address
	newAttrs.Address6 = address6
	newAttrs.CheckCommand = checkCommand
	newAttrs.Vars = variables
	newAttrs.Templates = templates
	newAttrs.Zone = zone

	if groups == nil {
		groups = []string{}
	}
	newAttrs.Groups = groups

	var newHost HostStruct
	newHost.Name = hostname
	newHost.Type = "Host"
	newHost.Attrs = newAttrs

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(newHost)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// Create the host
	results, err := server.NewAPIRequest(
		http.MethodPut,
		fmt.Sprintf("/objects/hosts/%s", url.PathEscape(hostname)),
		[]byte(payloadJSON),
	)

	// Ignore context deadline exceeded
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	}

	// Detect real errors
	if err == nil && results.Code != 200 {
		return nil, fmt.Errorf("%s", results.ErrorString)
	}

	// Wait for the host to be created
	operation := func() ([]HostStruct, error) {
		hosts, err := server.GetHost(hostname)
		if err != nil {
			return nil, backoff.Permanent(err)
		}
		for _, host := range hosts {
			if host.Name == hostname {
				return hosts, nil
			}
		}
		return nil, fmt.Errorf("host '%s' not found after creation", hostname)
	}

	// Number of tries must be at least 1 to avoid infinite loop
	tries := uint(math.Max(float64(server.Tries), 1.0))

	return backoff.Retry(
		context.Background(),
		operation,
		backoff.WithBackOff(backoff.NewConstantBackOff(server.RetryDelay)),
		backoff.WithMaxTries(tries),
	)
}

// UpdateHost updates a Host with its attrs
func (server *Server) UpdateHost(name string, attrs HostAttrs) ([]HostStruct, error) {

	host := HostStruct{
		Attrs: attrs,
	}

	body, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}

	r, err := server.NewAPIRequest(http.MethodPost, "/objects/hosts/"+name, body)
	if err != nil {
		return nil, err
	}

	if r.Code != http.StatusOK {
		return nil, fmt.Errorf("expected %d, got %d", http.StatusOK, r.Code)
	}

	jsonResponse, err := json.Marshal(r.Results)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the host response: %v", err)
	}

	var results []HostUpdateResult
	err = json.Unmarshal(jsonResponse, &results)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the host response: %v", err)
	}

	for _, result := range results {
		if result.Code != http.StatusOK {
			return nil, fmt.Errorf("%s", result.Status)
		}
	}

	return server.GetHost(name)
}

// DeleteHost deletes a host.
// When a context deadline is exceeded, wait for the host to be deleted if a number of tries is defined.
func (server *Server) DeleteHost(hostname string) error {
	results, err := server.NewAPIRequest(
		http.MethodDelete,
		fmt.Sprintf("/objects/hosts/%s?cascade=1", url.PathEscape(hostname)),
		nil,
	)

	// Ignore context deadline exceeded
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	// Detect real errors
	if err == nil && results.Code != 200 {
		return fmt.Errorf("%s", results.ErrorString)
	}

	// Wait for the host to be deleted
	operation := func() (string, error) {
		exists, err := server.HostExists(hostname)
		if err != nil {
			return "", backoff.Permanent(err)
		}
		if exists {
			return "", fmt.Errorf("host '%s' still exists after deletion", hostname)
		}
		return "", nil
	}

	// Number of tries must be at least 1 to avoid infinite loop
	tries := uint(math.Max(float64(server.Tries), 1.0))

	_, err = backoff.Retry(
		context.Background(),
		operation,
		backoff.WithBackOff(backoff.NewConstantBackOff(server.RetryDelay)),
		backoff.WithMaxTries(tries),
	)
	return err
}

// HostExists returns true if a Host exists
func (server *Server) HostExists(hostname string) (bool, error) {
	hosts, err := server.GetHost(hostname)
	if err != nil {
		return false, err
	}

	for _, host := range hosts {
		if host.Name == hostname {
			return true, nil
		}
	}

	return false, nil
}
