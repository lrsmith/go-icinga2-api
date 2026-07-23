package iapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const hostgroupEndpoint = "/objects/hostgroups"

// GetHostgroup fetches a HostGroup by its name.
func (server *Server) GetHostgroup(ctx context.Context, name string) ([]HostgroupStruct, error) {
	var hostgroups []HostgroupStruct
	endpoint := fmt.Sprintf("%v/%v", hostgroupEndpoint, name)
	_, err := server.NewAPIRequest(ctx, http.MethodGet, endpoint, nil, &hostgroups)
	if err != nil {
		return nil, err
	}

	if len(hostgroups) == 0 {
		return nil, nil
	}

	if len(hostgroups) != 1 {
		return nil, errors.New("found more than one matching hostgroup")
	}

	return hostgroups, nil
}

// CreateHostgroup creates a new HostGroup with its name and display name.
func (server *Server) CreateHostgroup(ctx context.Context, name, displayName, zone string) ([]HostgroupStruct, error) {
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
	results, err := server.NewAPIRequest(ctx, http.MethodPut, endpoint, payloadJSON, nil)
	if err != nil {
		return nil, err
	}

	if results.Code == http.StatusOK {
		return server.GetHostgroup(ctx, name)
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
}

// UpdateHostgroup updates a HostGroup with its attrs.
func (server *Server) UpdateHostgroup(ctx context.Context, name string, attrs HostgroupAttrs) ([]HostgroupStruct, error) {
	var hostgroup HostgroupStruct
	hostgroup.Attrs = attrs

	body, err := json.Marshal(hostgroup)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%v/%v", hostgroupEndpoint, name)
	r, err := server.NewAPIRequest(ctx, http.MethodPost, endpoint, body, nil)
	if err != nil {
		return nil, err
	}

	if r.Code != http.StatusOK {
		return nil, fmt.Errorf("expected %d, got %d", http.StatusOK, r.Code)
	}

	var results []HostgroupUpdateResult
	if len(r.Results) > 0 {
		if unmarshalErr := json.Unmarshal(r.Results, &results); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to unmarshal the host group response: %v", unmarshalErr)
		}
	}

	for _, result := range results {
		if result.Code != http.StatusOK {
			return nil, fmt.Errorf("%s", result.Status)
		}
	}

	return server.GetHostgroup(ctx, name)
}

// DeleteHostgroup deletes a HostGroup by its name.
func (server *Server) DeleteHostgroup(ctx context.Context, name string) error {
	endpoint := fmt.Sprintf("%v/%v", hostgroupEndpoint, name)
	results, err := server.NewAPIRequest(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return err
	}

	if results.Code == http.StatusOK {
		return nil
	}

	return fmt.Errorf("%s", results.ErrorString)
}

// HostgroupExists returns true if a HostGroup exists
func (server *Server) HostgroupExists(ctx context.Context, name string) (bool, error) {
	hostgroups, err := server.GetHostgroup(ctx, name)
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
