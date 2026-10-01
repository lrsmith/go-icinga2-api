package iapi

import (
	"context"
	"slices"
	"strings"
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

	_, err := Icinga2_Server.CreateCheckcommand(context.Background(), name, command, nil, nil)

	if err != nil {
		t.Error(err)
	}

}

func TestCreateCheckcommandWithTemplates(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	name := "check-command-templates"
	command := "/dev/null"
	commandArgs := map[string]string{"-X": "Xarg"}
	templates := []string{"go-icinga2-api-test-command"}

	checkcommands, err := Icinga2_Server.CreateCheckcommand(context.Background(), name, command, commandArgs, templates)
	if err != nil {
		t.Fatal(err)
	}

	for _, checkcommand := range checkcommands {
		if checkcommand.Name != name {
			continue
		}

		if !slices.Contains(checkcommand.Attrs.Templates, templates[0]) {
			t.Errorf("expected %s to import %s, got %v", name, templates[0], checkcommand.Attrs.Templates)
		}

		arguments, ok := checkcommand.Attrs.Arguments.(map[string]interface{})
		if !ok {
			t.Fatalf("expected arguments on %s, got %v", name, checkcommand.Attrs.Arguments)
		}
		if arguments["-X"] != "Xarg" {
			t.Errorf("expected argument -X to be Xarg, got %v", arguments["-X"])
		}
		// An argument of the template only reaches the checkcommand when the
		// template is imported instead of being stored as an attribute
		if arguments["-H"] != "$host.address$" {
			t.Errorf("expected argument -H to be inherited from the template, got %v", arguments["-H"])
		}
	}

	err = Icinga2_Server.DeleteCheckcommand(context.Background(), name)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateCheckcommandWithUnknownTemplate(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	name := "check-command-unknown-template"
	command := "/dev/null"
	templates := []string{"go-icinga2-api-test-does-not-exist"}

	_, err := Icinga2_Server.CreateCheckcommand(context.Background(), name, command, nil, templates)
	if err == nil {
		t.Fatalf("expected the creation of %s to be rejected", name)
	}

	if !strings.Contains(err.Error(), "Import references unknown template") {
		t.Errorf("expected the reason of the rejection, got %s", err)
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

// The templates of a checkcommand are immutable, Icinga reports a per object
// error while answering 200
func TestUpdateCheckcommandWithTemplates(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	name := "check-command-docker"
	attrs := CheckcommandAttrs{
		Command:   []string{"/bin/true"},
		Templates: []string{"plugin-check-command"},
	}

	_, err := Icinga2_Server.UpdateCheckcommand(context.Background(), name, attrs)
	if err == nil {
		t.Fatal("expected the update to be rejected")
	}

	if !strings.Contains(err.Error(), "Attribute cannot be modified") {
		t.Errorf("expected the reason of the rejection, got %s", err)
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

	_, err := Icinga2_Server.CreateCheckcommand(context.Background(), name, command, commandArgs, nil)
	if err != nil {
		t.Error(err)
	}

	// Delete check command after creating it.
	err = Icinga2_Server.DeleteCheckcommand(context.Background(), name)
	if err != nil {
		t.Error(err)
	}

}
