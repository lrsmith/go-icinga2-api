package iapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetService ...
func (server *Server) GetService(ctx context.Context, servicename, hostname string) ([]ServiceStruct, error) {
	var services []ServiceStruct
	_, err := server.NewAPIRequest(ctx, "GET", "/objects/services/"+hostname+"!"+servicename, nil, &services)
	if err != nil {
		return nil, err
	}
	return services, nil
}

// CreateService ...
func (server *Server) CreateService(ctx context.Context, servicename, hostname, checkCommand string, variables map[string]string, templates []string) ([]ServiceStruct, error) {
	var newAttrs ServiceAttrs
	newAttrs.CheckCommand = checkCommand
	newAttrs.Vars = variables
	newAttrs.Templates = templates

	var newService ServiceStruct
	newService.Attrs = newAttrs

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(newService)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// Make the API request to create the hosts.
	results, err := server.NewAPIRequest(ctx, "PUT", "/objects/services/"+hostname+"!"+servicename, payloadJSON, nil)
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		return server.GetService(ctx, servicename, hostname)
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
}

// UpdateService updates a Service with its attrs in-place
func (server *Server) UpdateService(ctx context.Context, servicename, hostname string, attrs ServiceAttrs) ([]ServiceStruct, error) {
	service := ServiceStruct{
		Attrs: attrs,
	}

	body, err := json.Marshal(service)
	if err != nil {
		return nil, err
	}

	r, err := server.NewAPIRequest(ctx, "POST", "/objects/services/"+hostname+"!"+servicename, body, nil)
	if err != nil {
		return nil, err
	}

	if r.Code != http.StatusOK {
		return nil, fmt.Errorf("expected %d, got %d", http.StatusOK, r.Code)
	}

	return server.GetService(ctx, servicename, hostname)
}

// DeleteService ...
func (server *Server) DeleteService(ctx context.Context, servicename, hostname string) error {
	results, err := server.NewAPIRequest(ctx, "DELETE", "/objects/services/"+hostname+"!"+servicename+"?cascade=1", nil, nil)
	if err != nil {
		return err
	}

	if results.Code == 200 {
		return nil
	}

	return fmt.Errorf("%s", results.ErrorString)
}
