package iapi

import (
	"encoding/json"
	"errors"
)

// GetHost ...
func (server *Server) GetHost(hostname string) ([]HostStruct, error) {

	var hosts []HostStruct

	results, err := server.NewAPIRequest("GET", "/objects/hosts/"+hostname, nil)
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

	if unmarshalErr := json.Unmarshal(jsonStr, &hosts); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return hosts, err

}

// CreateHost ...
func (server *Server) CreateHost(hostname, address, checkCommand string, variables map[string]string) ([]HostStruct, error) {

	var newAttrs HostAttrs
	newAttrs.Address = address
	newAttrs.CheckCommand = "hostalive"
	newAttrs.Vars = variables

	var newHost HostStruct
	newHost.Name = hostname
	newHost.Type = "Host"
	newHost.Attrs = newAttrs

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(newHost)
	if marshalErr != nil {
		return nil, marshalErr
	}

	//fmt.Printf("<payload> %s\n", payloadJSON)

	// Make the API request to create the hosts.
	results, err := server.NewAPIRequest("PUT", "/objects/hosts/"+hostname, []byte(payloadJSON))
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		hosts, err := server.GetHost(hostname)
		return hosts, err
	}

	// TODO Parse results.Results to get error messag
	return nil, errors.New(results.Status)

}

// DeleteHost ...
func (server *Server) DeleteHost(hostname string) error {

	results, err := server.NewAPIRequest("DELETE", "/objects/hosts/"+hostname+"?cascade=1", nil)
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
