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

// serviceCreateRequest is the payload of a service creation. Attributes are
// addressed by path, which is what makes "vars.os" merge with the variables
// inherited from the templates where a whole "vars" dictionary would replace them.
type serviceCreateRequest struct {
	Attrs     map[string]interface{} `json:"attrs"`
	Templates []string               `json:"templates,omitempty"`
}

// CreateService creates a service.
// The keys of variables are variable names, without the "vars." prefix.
func (server *Server) CreateService(ctx context.Context, servicename, hostname, checkCommand string, variables map[string]string, templates []string) ([]ServiceStruct, error) {
	// Only send the attributes the caller provided, an empty one would override
	// the value inherited from the templates.
	attrs := make(map[string]interface{})
	if checkCommand != "" {
		attrs["check_command"] = checkCommand
	}

	// Addressing each variable merges them with the ones inherited from the
	// templates, assigning the whole "vars" dictionary would replace them.
	for name, value := range variables {
		attrs["vars."+name] = value
	}

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(serviceCreateRequest{
		Attrs: attrs,
		// Templates are imports rather than an attribute. Icinga stores an
		// "attrs.templates" array on the object verbatim and imports nothing.
		Templates: templates,
	})
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

	return nil, objectCreateError(results)
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

	if updateErr := objectUpdateError(r); updateErr != nil {
		return nil, updateErr
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
