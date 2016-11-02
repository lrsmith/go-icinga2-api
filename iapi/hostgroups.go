package iapi

import (
	"encoding/json"
	"errors"
)

// GetHostgroup ...
func (server *Server) GetHostgroup(name string) ([]HostgroupStruct, error) {

	results, err := server.NewAPIRequest("GET", "/objects/hostgroups/"+name, nil)
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

	var hostgroup []HostgroupStruct
	if unmarshalErr := json.Unmarshal(jsonStr, &hostgroup); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return hostgroup, err

}

// CreateHostgroup ...
func (server *Server) CreateHostgroup(name, displayName string) error {

	var newAttrs HostgroupAttrs
	newAttrs.DisplayName = displayName

	var newHostgroup HostgroupStruct
	newHostgroup.Name = name
	newHostgroup.Type = "Hostgroup"
	newHostgroup.Attrs = newAttrs

	payloadJSON, marshalErr := json.Marshal(newHostgroup)
	if marshalErr != nil {
		return marshalErr
	}

	results, err := server.NewAPIRequest("PUT", "/objects/hostgroups/"+name, []byte(payloadJSON))
	if err != nil {
		return err
	}

	if results.Code == 200 {
		return nil
	}
	// TODO Parse results.Results to get error messag
	return errors.New(results.Status)

}

// DeleteHostgroup ...
func (server *Server) DeleteHostgroup(name string) error {

	results, err := server.NewAPIRequest("DELETE", "/objects/hostgroups/"+name, nil)
	if err != nil {
		return err
	}

	if results.Code == 200 {
		return nil
	} else if results.Code == 404 {
		if results.Status == "No objects found." {
			return nil
		}

	} else {
		return errors.New(results.Status)
	}

	return errors.New(results.Status)

}
