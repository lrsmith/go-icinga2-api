package iapi

import (
	"context"
	"testing"
)

func TestGetValidCheckcommand(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	name := "apache-status"

	_, err := Icinga2_Server.GetCheckcommand(context.Background(), name)

	if err != nil {
		t.Error(err)
	}
}

func TestGetInvalidCheckcommand(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	name := "invalid-check-command"

	_, err := Icinga2_Server.GetCheckcommand(context.Background(), name)
	if err != nil {
		t.Error(err)
	}

}

func TestCreateCheckcommand(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	name := "check-command-docker"
	command := "/dev/null"

	_, err := Icinga2_Server.CreateCheckcommand(context.Background(), name, command, nil)

	if err != nil {
		t.Error(err)
	}

}

func TestUpdateCheckcommand(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	name := "check-command-docker"
	attrs := CheckcommandAttrs{
		Command: []string{"/bin/true"},
		Arguments: map[string]string{
			"-Y": "Yarg",
		},
	}

	_, err := Icinga2_Server.UpdateCheckcommand(context.Background(), name, attrs)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteCheckcommand(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	name := "check-command-docker"

	err := Icinga2_Server.DeleteCheckcommand(context.Background(), name)
	if err != nil {
		t.Error(err)
	}

}

func TestCreateCheckcommandArgs(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	name := "check-command-docker-args"
	command := "/dev/null"
	commandArgs := make(map[string]string)
	commandArgs["-I"] = "Iarg"
	commandArgs["-X"] = "Xarg"

	_, err := Icinga2_Server.CreateCheckcommand(context.Background(), name, command, commandArgs)
	if err != nil {
		t.Error(err)
	}

	// Delete check command after creating it.
	err = Icinga2_Server.DeleteCheckcommand(context.Background(), name)
	if err != nil {
		t.Error(err)
	}

}
