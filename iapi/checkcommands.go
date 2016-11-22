package iapi

import (
	"encoding/json"
	"errors"
)

// GetCheckCommand ...
func (server *Server) GetCheckCommand(name string) (*CheckCommandStruct, error) {

	results, err := server.NewAPIRequest("GET", "/objects/checkcommands/"+name, nil)
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

	var checkcommand []CheckCommandStruct
	if unmarshalErr := json.Unmarshal(jsonStr, &checkcommand); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	if len(checkcommand) == 0 {
		return nil, nil
	}
	if len(checkcommand) != 1 {
		return nil, errors.New("Found more than one matching check command.")
	}

	return &checkcommand[0], err

}

// CreateCheckCommand ...
func (server *Server) CreateCheckCommand(name, command string, command_arguments map[string]string) (*CheckCommandStruct, error) {

	var newAttrs CheckCommandAttrs
	newAttrs.Command = []string{command}
	newAttrs.Arguments = command_arguments

	var newCheckCommand CheckCommandStruct
	newCheckCommand.Name = name
	newCheckCommand.Type = "CheckCommand"
	newCheckCommand.Attrs = newAttrs

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(newCheckCommand)
	if marshalErr != nil {
		return nil, marshalErr
	}

	//fmt.Printf("<payload> %s\n", payloadJSON)

	// Make the API request to create the hosts.
	results, err := server.NewAPIRequest("PUT", "/objects/checkcommands/"+name, []byte(payloadJSON))
	if err != nil {
		return nil, err
	}

	//fmt.Printf("%v", results) // Useful debug. Better error message.
	if results.Code == 200 {
		theCheckCommand, err := server.GetCheckCommand(name)
		return theCheckCommand, err
	}

	// TODO Parse results.Results to get error messag
	return nil, errors.New(results.Status)

}

// DeleteCheckCommand ...
func (server *Server) DeleteCheckCommand(name string) error {

	results, err := server.NewAPIRequest("DELETE", "/objects/checkcommands/"+name+"?cascade=1", nil)
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
