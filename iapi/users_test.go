package iapi

import (
	"strings"
	"testing"
)

func TestUser(t *testing.T) {
	icingaServer := Server{"root", ICINGA2_API_PASSWORD, "https://127.0.0.1:5665/v1", true, nil}
	t.Run("Create", func(t *testing.T) {
		t.Run("SimpleUser", func(t *testing.T) {
			username := "test-user"
			_, err := icingaServer.CreateUser(username, "")
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("UserAlreadyExists", func(t *testing.T) {
			username := "test-user"
			_, err := icingaServer.CreateUser(username, "")
			if !strings.HasSuffix(err.Error(), "500 Object could not be created") {
				t.Error(err)
			}
		})

		t.Run("UserWithEmail", func(t *testing.T) {
			username := "test-user-with-email"
			email := "email@example.com"
			_, err := icingaServer.CreateUser(username, email)
			if err != nil {
				t.Error(err)
			}

			// Delete user after creating it.
			err = icingaServer.DeleteUser(username)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("Read", func(t *testing.T) {
		t.Run("ValidUser", func(t *testing.T) {
			username := "valid-user"
			_, err := icingaServer.GetUser(username)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("InvalidUser", func(t *testing.T) {
			username := "invalid-user"
			_, err := icingaServer.GetUser(username)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		username := "test-user"
		t.Run("UserExists", func(t *testing.T) {
			err := icingaServer.DeleteUser(username)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("UserDoNotExist", func(t *testing.T) {
			err := icingaServer.DeleteUser(username)
			if err.Error() != "No objects found." {
				t.Error(err)
			}
		})
	})
}
