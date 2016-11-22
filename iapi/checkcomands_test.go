package iapi

import "testing"

func TestGetValidCheckCommand(t *testing.T) {

	name := "apache-status"

	checkcommand, err := VagrantImage.GetCheckCommand(name)

	if (err != nil) || (checkcommand == nil) {
		t.Errorf("Error : Failed to find check command %s : ( %s <> %v ) ", name, err, checkcommand)
	}

	if checkcommand.Name != name {
		t.Errorf("Error : Did not get expected check command. ( %s != %s )", checkcommand.Name, name)
	}

}

func TestGetInvalidCheckCommand(t *testing.T) {

	name := "invalid-check-command"

	checkcommand, err := VagrantImage.GetCheckCommand(name)
	if err != nil && checkcommand != nil {
		t.Errorf("Error : Did not get empty list of check commands. ( %s : %v )", err, checkcommand)
	}

}

func TestCreateCheckCommand(t *testing.T) {

	name := "check-command-docker"
	command := "/dev/null"

	_, err := VagrantImage.CreateCheckCommand(name, command, nil)

	if err != nil {
		t.Errorf("Error : Failed to create Check Command %s : %s", name, err)
	}

}

func TestDeleteCheckCommand(t *testing.T) {

	name := "check-command-docker"

	err := VagrantImage.DeleteCheckCommand(name)
	if err != nil {
		t.Errorf("Error : Failed to delete Check Command %s : %s", name, err)
	}

}

func TestCreateCheckCommandArgs(t *testing.T) {

	name := "check-command-docker-args"
	command := "/dev/null"
	command_args := make(map[string]string)
	command_args["-I"] = "Iarg"
	command_args["-X"] = "Xarg"

	_, err := VagrantImage.CreateCheckCommand(name, command, command_args)

	if err != nil {
		t.Errorf("Error : Failed to create Check Command %s : %s", name, err)
	}

	// Delete check command after creating it.
	deleteErr := VagrantImage.DeleteCheckCommand(name)
	if deleteErr != nil {
		t.Errorf("Error : Cleanup failed for Check Command %s : %s", name, err)
	}

}

/* NOT WORKING
func TestDeleteNonExistentCheckCommand(t *testing.T) {

	name := "check-command-docker-2"
	err := VagrantImage.DeleteCheckCommand(name)
	if err != nil {
		t.Errorf("Error : Failed to delete Check Command %s : %s", name, err)
	}

}
*/
