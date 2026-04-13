package iapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const hostgroupEndpoint = "/objects/hostgroups"

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
func (server *Server) CreateHostgroup(name, displayName, zone string) ([]HostgroupStruct, error) {
	var newAttrs HostgroupAttrs
	newAttrs.DisplayName = displayName
	newAttrs.Zone = zone

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

// UpdateHostgroup updates a HostGroup with its attrs.
func (server *Server) UpdateHostgroup(name string, attrs HostgroupAttrs) ([]HostgroupStruct, error) {
	var hostgroup HostgroupStruct
	hostgroup.Attrs = attrs

	body, err := json.Marshal(hostgroup)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%v/%v", hostgroupEndpoint, name)
	r, err := server.NewAPIRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}

	if r.Code != http.StatusOK {
		return nil, fmt.Errorf("expected %d, got %d", http.StatusOK, r.Code)
	}

	jsonResponse, err := json.Marshal(r.Results)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the host group response: %v", err)
	}

	var results []HostgroupUpdateResult
	err = json.Unmarshal(jsonResponse, &results)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the host group response: %v", err)
	}

	for _, result := range results {
		if result.Code != http.StatusOK {
			return nil, fmt.Errorf("%s", result.Status)
		}
	}

	return server.GetHostgroup(name)
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

// HostgroupExists returns true if a HostGroup exists
func (server *Server) HostgroupExists(name string) (bool, error) {
	hostgroups, err := server.GetHostgroup(name)
	if err != nil {
		return false, err
	}

	for _, hostgroup := range hostgroups {
		if hostgroup.Name == name {
			return true, nil
		}
	}

	return false, nil
}
