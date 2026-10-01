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

// hostCreateRequest is the payload of a host creation. Attributes are addressed
// by path, which is what makes "vars.os" merge with the variables inherited from
// the templates where a whole "vars" dictionary would replace them.
type hostCreateRequest struct {
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Attrs     map[string]interface{} `json:"attrs"`
	Templates []string               `json:"templates,omitempty"`
}

// CreateHost creates a host.
// The keys of variables are variable names, without the "vars." prefix.
// When a context deadline is exceeded, wait for the host to be created if a number of tries is defined.
func (server *Server) CreateHost(ctx context.Context, hostname, address, address6 string, checkCommand string, variables map[string]interface{}, templates []string, groups []string, zone string) (hosts []HostStruct, err error) {
	// Only send the attributes the caller provided, an empty one would override
	// the value inherited from the templates.
	attrs := make(map[string]interface{})
	if address != "" {
		attrs["address"] = address
	}
	if address6 != "" {
		attrs["address6"] = address6
	}
	if checkCommand != "" {
		attrs["check_command"] = checkCommand
	}
	if zone != "" {
		attrs["zone"] = zone
	}
	if groups != nil {
		attrs["groups"] = groups
	}

	// Addressing each variable merges them with the ones inherited from the
	// templates, assigning the whole "vars" dictionary would replace them.
	for name, value := range variables {
		attrs["vars."+name] = value
	}

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(hostCreateRequest{
		Name:  hostname,
		Type:  "Host",
		Attrs: attrs,
		// Templates are imports rather than an attribute. Icinga stores an
		// "attrs.templates" array on the object verbatim and imports nothing.
		Templates: templates,
	})
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
		return nil, objectCreateError(results)
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

	if updateErr := objectUpdateError(r); updateErr != nil {
		return nil, updateErr
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
