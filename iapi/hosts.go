package iapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetHost ...
func (server *Server) GetHost(hostname string) ([]HostStruct, error) {

	var hosts []HostStruct

	results, err := server.NewAPIRequest(http.MethodGet, "/objects/hosts/"+hostname, nil)
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
func (server *Server) CreateHost(hostname, address, address6 string, checkCommand string, variables map[string]interface{}, templates []string, groups []string, zone string) ([]HostStruct, error) {

	var newAttrs HostAttrs
	newAttrs.Address = address
	newAttrs.Address6 = address6
	newAttrs.CheckCommand = checkCommand
	newAttrs.Vars = variables
	newAttrs.Templates = templates
	newAttrs.Zone = zone

	if groups == nil {
		groups = []string{}
	}
	newAttrs.Groups = groups

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
	results, err := server.NewAPIRequest(http.MethodPut, "/objects/hosts/"+hostname, []byte(payloadJSON))
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		hosts, err := server.GetHost(hostname)
		return hosts, err
	}

	return nil, fmt.Errorf("%s", results.ErrorString)

}

// UpdateHost updates a Host with its attrs
func (server *Server) UpdateHost(name string, attrs HostAttrs) ([]HostStruct, error) {

	host := HostStruct{
		Attrs: attrs,
	}

	body, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}

	r, err := server.NewAPIRequest(http.MethodPost, "/objects/hosts/"+name, body)
	if err != nil {
		return nil, err
	}

	if r.Code != http.StatusOK {
		return nil, fmt.Errorf("expected %d, got %d", http.StatusOK, r.Code)
	}

	jsonResponse, err := json.Marshal(r.Results)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the host response: %v", err)
	}

	var results []HostUpdateResult
	err = json.Unmarshal(jsonResponse, &results)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the host response: %v", err)
	}

	for _, result := range results {
		if result.Code != http.StatusOK {
			return nil, fmt.Errorf("%s", result.Status)
		}
	}

	return server.GetHost(name)
}

// DeleteHost ...
func (server *Server) DeleteHost(hostname string) error {
	results, err := server.NewAPIRequest(http.MethodDelete, "/objects/hosts/"+hostname+"?cascade=1", nil)
	if err != nil {
		return err
	}

	if results.Code == 200 {
		return nil
	}

	return fmt.Errorf("%s", results.ErrorString)
}

// HostExists returns true if a Host exists
func (server *Server) HostExists(hostname string) (bool, error) {
	hosts, err := server.GetHost(hostname)
	if err != nil {
		return false, err
	}

	for _, host := range hosts {
		if host.Name == hostname {
			return true, nil
		}
	}

	return false, nil
}
