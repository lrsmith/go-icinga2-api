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

// CreateCheckcommand ...
func (server *Server) CreateCheckcommand(ctx context.Context, name, command string, commandArguments map[string]string) ([]CheckcommandStruct, error) {
	var newAttrs CheckcommandAttrs
	newAttrs.Command = []string{command}
	newAttrs.Arguments = commandArguments

	var newCheckcommand CheckcommandStruct
	newCheckcommand.Name = name
	newCheckcommand.Type = "CheckCommand"
	newCheckcommand.Attrs = newAttrs

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(newCheckcommand)
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

	return nil, fmt.Errorf("%s", results.ErrorString)
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
