package iapi

import "testing"

func TestGetValidCheckcommand(t *testing.T) {

	name := "apache-status"

	_, err := Icinga2_Server.GetCheckcommand(name)

	if err != nil {
		t.Error(err)
	}
}

func TestGetInvalidCheckcommand(t *testing.T) {

	name := "invalid-check-command"

	_, err := Icinga2_Server.GetCheckcommand(name)
	if err != nil {
		t.Error(err)
	}

}

func TestCreateCheckcommand(t *testing.T) {

	name := "check-command-docker"
	command := "/dev/null"

	_, err := Icinga2_Server.CreateCheckcommand(name, command, nil)

	if err != nil {
		t.Error(err)
	}

}

func TestDeleteCheckcommand(t *testing.T) {

	name := "check-command-docker"

	err := Icinga2_Server.DeleteCheckcommand(name)
	if err != nil {
		t.Error(err)
	}

}

func TestCreateCheckcommandArgs(t *testing.T) {
	name := "check-command-docker-args"
	command := "/dev/null"
	commandArgs := make(map[string]string)
	commandArgs["-I"] = "Iarg"
	commandArgs["-X"] = "Xarg"

	_, err := Icinga2_Server.CreateCheckcommand(name, command, commandArgs)
	if err != nil {
		t.Error(err)
	}

	// Delete check command after creating it.
	err = Icinga2_Server.DeleteCheckcommand(name)
	if err != nil {
		t.Error(err)
	}

}
