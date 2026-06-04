package iapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const apiUserEndpoint = "/objects/apiusers"

// GetApiUser fetches an ApiUser by its name.
func (server *Server) GetApiUser(name string) ([]ApiUserStruct, error) {
	endpoint := fmt.Sprintf("%v/%v", apiUserEndpoint, name)
	results, err := server.NewAPIRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Contents of the results is an interface object. Need to convert it to json first.
	jsonStr, err := json.Marshal(results.Results)
	if err != nil {
		return nil, err
	}

	// then the JSON can be pushed into the appropriate struct.
	// Note : Results is a slice so much push into a slice.
	var apiUsers []ApiUserStruct
	if err := json.Unmarshal(jsonStr, &apiUsers); err != nil {
		return nil, err
	}

	if len(apiUsers) == 0 {
		return nil, nil
	}

	if len(apiUsers) != 1 {
		return nil, errors.New("found more than one matching apiuser")
	}

	return apiUsers, err
}

// CreateApiUser creates a new ApiUser with its name, password, clientCNs, and permissions.
func (server *Server) CreateApiUser(name, password string, clientCN string, permissions []string) ([]ApiUserStruct, error) {
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
	results, err := server.NewAPIRequest(http.MethodPut, endpoint, payloadJSON)
	if err != nil {
		return nil, err
	}

	if results.Code == http.StatusOK {
		apiUsers, err := server.GetApiUser(name)
		return apiUsers, err
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
}

// UpdateApiUser updates an ApiUser with its params.
func (server *Server) UpdateApiUser(name string, params *ApiUserAttrs) ([]ApiUserStruct, error) {
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
	results, err := server.NewAPIRequest(http.MethodPost, endpoint, attrsBody)
	if err != nil {
		return nil, err
	}

	if results.Code == http.StatusOK {
		apiUsers, err := server.GetApiUser(name)
		return apiUsers, err
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
}

// DeleteApiUser deletes an ApiUser by its name.
func (server *Server) DeleteApiUser(name string) error {
	endpoint := fmt.Sprintf("%v/%v", apiUserEndpoint, name)
	results, err := server.NewAPIRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	if results.Code == http.StatusOK {
		return nil
	}

	return fmt.Errorf("%s", results.ErrorString)
}

// ApiUserExists returns true if an ApiUser exists
func (server *Server) ApiUserExists(name string) (bool, error) {
	apiUsers, err := server.GetApiUser(name)
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
