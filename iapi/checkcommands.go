package iapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetCheckcommand ...
func (server *Server) GetCheckcommand(ctx context.Context, name string) ([]CheckcommandStruct, error) {
	var checkcommands []CheckcommandStruct
	_, err := server.NewAPIRequest(ctx, "GET", "/objects/checkcommands/"+name, nil, &checkcommands)
	if err != nil {
		return nil, err
	}
	return checkcommands, nil
}

// checkcommandCreateRequest is the payload of a checkcommand creation.
// Attributes are addressed by path, which is what makes "arguments.-I" merge
// with the arguments inherited from the templates where a whole "arguments"
// dictionary would replace them.
type checkcommandCreateRequest struct {
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Attrs     map[string]interface{} `json:"attrs"`
	Templates []string               `json:"templates,omitempty"`
}

// CreateCheckcommand creates a checkcommand.
func (server *Server) CreateCheckcommand(ctx context.Context, name, command string, commandArguments map[string]string, templates []string) ([]CheckcommandStruct, error) {
	// Only send the attributes the caller provided, an empty one would override
	// the value inherited from the templates.
	attrs := make(map[string]interface{})
	if command != "" {
		attrs["command"] = []string{command}
	}

	// Addressing each argument merges them with the ones inherited from the
	// templates, assigning the whole "arguments" dictionary would replace them.
	for argument, value := range commandArguments {
		attrs["arguments."+argument] = value
	}

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(checkcommandCreateRequest{
		Name:  name,
		Type:  "CheckCommand",
		Attrs: attrs,
		// Templates are imports rather than an attribute. Icinga stores an
		// "attrs.templates" array on the object verbatim and imports nothing.
		Templates: templates,
	})
	if marshalErr != nil {
		return nil, marshalErr
	}

	// Make the API request to create the checkcommands.
	results, err := server.NewAPIRequest(ctx, "PUT", "/objects/checkcommands/"+name, payloadJSON, nil)
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		return server.GetCheckcommand(ctx, name)
	}

	return nil, objectCreateError(results)
}

// UpdateCheckcommand updates a CheckCommand with its attrs in-place
func (server *Server) UpdateCheckcommand(ctx context.Context, name string, attrs CheckcommandAttrs) ([]CheckcommandStruct, error) {
	checkcommand := CheckcommandStruct{
		Attrs: attrs,
	}

	body, err := json.Marshal(checkcommand)
	if err != nil {
		return nil, err
	}

	r, err := server.NewAPIRequest(ctx, "POST", "/objects/checkcommands/"+name, body, nil)
	if err != nil {
		return nil, err
	}

	if r.Code != http.StatusOK {
		return nil, fmt.Errorf("expected %d, got %d", http.StatusOK, r.Code)
	}

	if updateErr := objectUpdateError(r); updateErr != nil {
		return nil, updateErr
	}

	return server.GetCheckcommand(ctx, name)
}

// DeleteCheckcommand ...
func (server *Server) DeleteCheckcommand(ctx context.Context, name string) error {
	results, err := server.NewAPIRequest(ctx, "DELETE", "/objects/checkcommands/"+name+"?cascade=1", nil, nil)
	if err != nil {
		return err
	}

	if results.Code == 200 {
		return nil
	}

	return fmt.Errorf("%s", results.ErrorString)
}
