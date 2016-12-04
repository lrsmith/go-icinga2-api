package iapi

import (
	"encoding/json"
	"errors"
)

// GetCheckCommand ...
func (server *Server) GetCheckCommand(name string) ([]CheckCommandStruct, error) {

	var checkcommands []CheckCommandStruct
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

	if unmarshalErr := json.Unmarshal(jsonStr, &checkcommands); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return checkcommands, err

}

// CreateCheckCommand ...
func (server *Server) CreateCheckCommand(name, command string, command_arguments map[string]string) ([]CheckCommandStruct, error) {

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
	// Need to check for 500 'already exsts as this is not necessarialy an error. Would require a deeper inspection to determine.
/*&{0 500 Object could not be created 500 [map[errors:[Configuration file '/var/lib/icinga2/api/packages/_api/icinga2-1476058673-1/conf.d/checkcommands/terraform-test-checkcommand.conf' already exists.] status:Object could not be created. code:500]]}--- FAIL: TestAccCreateCheckCommand (0.02s)*/

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
