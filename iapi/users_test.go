package iapi

import (
	"strings"
	"testing"
)

func TestGetValidUser(t *testing.T) {

	username := "valid-user"

	_, err := Icinga2_Server.GetUser(username)

	if err != nil {
		t.Error(err)
	}
}

func TestGetInvalidUser(t *testing.T) {

	username := "invalid-user"
	_, err := Icinga2_Server.GetUser(username)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateSimpleUser(t *testing.T) {

	username := "test-user"

	_, err := Icinga2_Server.CreateUser(username, "")

	if err != nil {
		t.Error(err)
	}
}

func TestCreateUserAlreadyExists(t *testing.T) {

	username := "test-user"

	_, err := Icinga2_Server.CreateUser(username, "")

	if !strings.HasSuffix(err.Error(), " already exists.") {
		t.Error(err)
	}
}

func TestCreateUserWithEmail(t *testing.T) {

	username := "test-user-with-email"
	Email := "email@example.com"

	_, err := Icinga2_Server.CreateUser(username, Email)
	if err != nil {
		t.Error(err)
	}

	// Delete user after creating it.
	deleteErr := Icinga2_Server.DeleteUser(username)
	if deleteErr != nil {
		t.Error(err)
	}
}

func TestDeleteUser(t *testing.T) {

	username := "test-user"

	err := Icinga2_Server.DeleteUser(username)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteUserDNE(t *testing.T) {
	username := "test-user"
	err := Icinga2_Server.DeleteUser(username)
	if err.Error() != "No objects found." {
		t.Error(err)
	}
}