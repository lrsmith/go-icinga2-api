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
func (server *Server) GetHost(ctx context.Context, hostname string) ([]HostStruct, error) {
	var hosts []HostStruct

	_, err := server.NewAPIRequest(ctx, http.MethodGet, "/objects/hosts/"+hostname, nil, &hosts)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

// CreateHost creates a host.
// When a context deadline is exceeded, wait for the host to be created if a number of tries is defined.
func (server *Server) CreateHost(ctx context.Context, hostname, address, address6 string, checkCommand string, variables map[string]interface{}, templates []string, groups []string, zone string) (hosts []HostStruct, err error) {
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
		ctx,
		http.MethodPut,
		fmt.Sprintf("/objects/hosts/%s", url.PathEscape(hostname)),
		payloadJSON,
		nil,
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
		hosts, err := server.GetHost(ctx, hostname)
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
		ctx,
		operation,
		backoff.WithBackOff(backoff.NewConstantBackOff(server.RetryDelay)),
		backoff.WithMaxTries(tries),
	)
}

// UpdateHost updates a Host with its attrs
func (server *Server) UpdateHost(ctx context.Context, name string, attrs HostAttrs) ([]HostStruct, error) {
	host := HostStruct{
		Attrs: attrs,
	}

	body, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}

	r, err := server.NewAPIRequest(ctx, http.MethodPost, "/objects/hosts/"+name, body, nil)
	if err != nil {
		return nil, err
	}

	if r.Code != http.StatusOK {
		return nil, fmt.Errorf("expected %d, got %d", http.StatusOK, r.Code)
	}

	var results []HostUpdateResult
	if len(r.Results) > 0 {
		if unmarshalErr := json.Unmarshal(r.Results, &results); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to unmarshal the host response: %v", unmarshalErr)
		}
	}

	for _, result := range results {
		if result.Code != http.StatusOK {
			return nil, fmt.Errorf("%s", result.Status)
		}
	}

	return server.GetHost(ctx, name)
}

// DeleteHost deletes a host.
// When a context deadline is exceeded, wait for the host to be deleted if a number of tries is defined.
func (server *Server) DeleteHost(ctx context.Context, hostname string) error {
	results, err := server.NewAPIRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("/objects/hosts/%s?cascade=1", url.PathEscape(hostname)),
		nil,
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
		exists, err := server.HostExists(ctx, hostname)
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
		ctx,
		operation,
		backoff.WithBackOff(backoff.NewConstantBackOff(server.RetryDelay)),
		backoff.WithMaxTries(tries),
	)
	return err
}

// HostExists returns true if a Host exists
func (server *Server) HostExists(ctx context.Context, hostname string) (bool, error) {
	hosts, err := server.GetHost(ctx, hostname)
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
