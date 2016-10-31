package icinga2

import "encoding/json"

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
