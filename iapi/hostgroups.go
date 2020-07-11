package iapi

import (
	"encoding/json"
	"errors"
	"fmt"
)

// GetHostgroup ...
func (server *Server) GetHostgroup(name string) ([]HostgroupStruct, error) {

// HostgroupParams defines all available options related to updating a HostGroup.
type HostgroupParams struct {
	DisplayName string
}
	if err != nil {
		return nil, err
	}

	// Contents of the results is an interface object. Need to convert it to json first.
	jsonStr, marshalErr := json.Marshal(results.Results)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// then the JSON can be pushed into the appropriate struct.
	// Note : Results is a slice so much push into a slice.

	if unmarshalErr := json.Unmarshal(jsonStr, &hostgroups); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	if len(hostgroups) == 0 {
		return nil, nil
	}

	if len(hostgroups) != 1 {
		return nil, errors.New("found more than one matching hostgroup")
	}

	return hostgroups, err

}

// CreateHostgroup ...
func (server *Server) CreateHostgroup(name, displayName string) ([]HostgroupStruct, error) {

	var newAttrs HostgroupAttrs
	newAttrs.DisplayName = displayName

	var newHostgroup HostgroupStruct
	newHostgroup.Name = name
	newHostgroup.Type = "Hostgroup"
	newHostgroup.Attrs = newAttrs

	payloadJSON, marshalErr := json.Marshal(newHostgroup)
	if marshalErr != nil {
		return nil, marshalErr
	}

	results, err := server.NewAPIRequest("PUT", "/objects/hostgroups/"+name, []byte(payloadJSON))
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		hostgroups, err := server.GetHostgroup(name)
		return hostgroups, err
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
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

// DeleteHostgroup ...
func (server *Server) DeleteHostgroup(name string) error {
	results, err := server.NewAPIRequest("DELETE", "/objects/hostgroups/"+name, nil)
	if err != nil {
		return err
	}

	if results.Code == 200 {
		return nil
	}

	return fmt.Errorf("%s", results.ErrorString)
}
