package iapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const hostgroupEndpoint = "/objects/hostgroups"

// HostgroupParams defines all available options related to updating a HostGroup.
type HostgroupParams struct {
	DisplayName string
}

// GetHostgroup fetches a HostGroup by its name.
func (server *Server) GetHostgroup(name string) ([]HostgroupStruct, error) {
	endpoint := fmt.Sprintf("%v/%v", hostgroupEndpoint, name)
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
	var hostgroups []HostgroupStruct
	if err := json.Unmarshal(jsonStr, &hostgroups); err != nil {
		return nil, err
	}

	if len(hostgroups) == 0 {
		return nil, nil
	}

	if len(hostgroups) != 1 {
		return nil, errors.New("found more than one matching hostgroup")
	}

	return hostgroups, err
}

// CreateHostgroup creates a new HostGroup with its name and display name.
func (server *Server) CreateHostgroup(name, displayName string) ([]HostgroupStruct, error) {
	var newAttrs HostgroupAttrs
	newAttrs.DisplayName = displayName

	var newHostgroup HostgroupStruct
	newHostgroup.Name = name
	newHostgroup.Type = "Hostgroup"
	newHostgroup.Attrs = newAttrs

	payloadJSON, err := json.Marshal(newHostgroup)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%v/%v", hostgroupEndpoint, name)
	results, err := server.NewAPIRequest(http.MethodPut, endpoint, payloadJSON)
	if err != nil {
		return nil, err
	}

	if results.Code == http.StatusOK {
		hostgroups, err := server.GetHostgroup(name)
		return hostgroups, err
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
}

// UpdateHostgroup updates a HostGroup with its params.
func (server *Server) UpdateHostgroup(name string, params *HostgroupParams) ([]HostgroupStruct, error) {
	attrs := make(map[string]interface{})
	if params.DisplayName != "" {
		attrs["display_name"] = params.DisplayName
	}

	attrsMap := map[string]interface{}{
		"attrs": attrs,
	}

	attrsBody, err := json.Marshal(attrsMap)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%v/%v", hostgroupEndpoint, name)
	results, err := server.NewAPIRequest(http.MethodPost, endpoint, attrsBody)
	if err != nil {
		return nil, err
	}

	if results.Code == http.StatusOK {
		hostgroups, err := server.GetHostgroup(name)
		return hostgroups, err
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
}

// DeleteHostgroup deletes a HostGroup by its name.
func (server *Server) DeleteHostgroup(name string) error {
	endpoint := fmt.Sprintf("%v/%v", hostgroupEndpoint, name)
	results, err := server.NewAPIRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	if results.Code == http.StatusOK {
		return nil
	}

	return fmt.Errorf("%s", results.ErrorString)
}
