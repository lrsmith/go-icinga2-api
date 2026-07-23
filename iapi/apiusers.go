package iapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const apiUserEndpoint = "/objects/apiusers"

// GetApiUser fetches an ApiUser by its name.
func (server *Server) GetApiUser(ctx context.Context, name string) ([]ApiUserStruct, error) {
	var apiUsers []ApiUserStruct
	endpoint := fmt.Sprintf("%v/%v", apiUserEndpoint, name)
	_, err := server.NewAPIRequest(ctx, http.MethodGet, endpoint, nil, &apiUsers)
	if err != nil {
		return nil, err
	}

	if len(apiUsers) == 0 {
		return nil, nil
	}

	if len(apiUsers) != 1 {
		return nil, errors.New("found more than one matching apiuser")
	}

	return apiUsers, nil
}

// CreateApiUser creates a new ApiUser with its name, password, clientCNs, and permissions.
func (server *Server) CreateApiUser(ctx context.Context, name, password string, clientCN string, permissions []string) ([]ApiUserStruct, error) {
	var newAttrs ApiUserAttrs
	newAttrs.Password = password
	newAttrs.ClientCN = clientCN

	if permissions == nil {
		permissions = []string{}
	}
	newAttrs.Permissions = permissions

	var newApiUser ApiUserStruct
	newApiUser.Name = name
	newApiUser.Type = "ApiUser"
	newApiUser.Attrs = newAttrs

	payloadJSON, err := json.Marshal(newApiUser)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%v/%v", apiUserEndpoint, name)
	results, err := server.NewAPIRequest(ctx, http.MethodPut, endpoint, payloadJSON, nil)
	if err != nil {
		return nil, err
	}

	if results.Code == http.StatusOK {
		return server.GetApiUser(ctx, name)
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
}

// UpdateApiUser updates an ApiUser with its params.
func (server *Server) UpdateApiUser(ctx context.Context, name string, params *ApiUserAttrs) ([]ApiUserStruct, error) {
	attrs := make(map[string]interface{})
	if params.Password != "" {
		attrs["password"] = params.Password
	}
	if params.Permissions != nil {
		attrs["permissions"] = params.Permissions
	}

	attrsMap := map[string]interface{}{
		"attrs": attrs,
	}

	attrsBody, err := json.Marshal(attrsMap)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%v/%v", apiUserEndpoint, name)
	results, err := server.NewAPIRequest(ctx, http.MethodPost, endpoint, attrsBody, nil)
	if err != nil {
		return nil, err
	}

	if results.Code == http.StatusOK {
		return server.GetApiUser(ctx, name)
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
}

// DeleteApiUser deletes an ApiUser by its name.
func (server *Server) DeleteApiUser(ctx context.Context, name string) error {
	endpoint := fmt.Sprintf("%v/%v", apiUserEndpoint, name)
	results, err := server.NewAPIRequest(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return err
	}

	if results.Code == http.StatusOK {
		return nil
	}

	return fmt.Errorf("%s", results.ErrorString)
}

// ApiUserExists returns true if an ApiUser exists
func (server *Server) ApiUserExists(ctx context.Context, name string) (bool, error) {
	apiUsers, err := server.GetApiUser(ctx, name)
	if err != nil {
		return false, err
	}

	for _, apiUser := range apiUsers {
		if apiUser.Name == name {
			return true, nil
		}
	}

	return false, nil
}
