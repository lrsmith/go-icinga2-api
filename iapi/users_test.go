package iapi

import (
	"context"
	"testing"
)

func TestUser(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}
	icingaServer, err := New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, ICINGA2_API_URL, ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		t.Run("SimpleUser", func(t *testing.T) {
			username := "test-user"
			_, err := icingaServer.CreateUser(ctx, username, "", nil)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("UserAlreadyExists", func(t *testing.T) {
			username := "test-user"
			_, err := icingaServer.CreateUser(ctx, username, "", nil)
			if err == nil {
				t.Error("expected error user already exists, got nil")
			}
		})

		t.Run("UserWithEmail", func(t *testing.T) {
			username := "test-user-with-email"
			email := "email@example.com"
			_, err := icingaServer.CreateUser(ctx, username, email, nil)
			if err != nil {
				t.Error(err)
			}

			// Delete user after creating it.
			err = icingaServer.DeleteUser(ctx, username)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("UserWithVars", func(t *testing.T) {
			username := "test-user-with-vars"

			variables := make(map[string]string)

			variables["vars.login"] = "test_user"
			variables["vars.ou"] = "custom support"

			_, err := icingaServer.CreateUser(ctx, username, "", variables)
			if err != nil {
				t.Error(err)
			}

			// Delete user after creating it.
			err = icingaServer.DeleteUser(ctx, username)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("Read", func(t *testing.T) {
		t.Run("ValidUser", func(t *testing.T) {
			username := "valid-user"
			_, err := icingaServer.GetUser(ctx, username)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("InvalidUser", func(t *testing.T) {
			username := "invalid-user"
			_, err := icingaServer.GetUser(ctx, username)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("ValidUser", func(t *testing.T) {
			username := "test-user"
			attrs := UserAttrs{
				Email: "test-updated@example.com",
				Vars: map[string]string{
					"vars.ou": "updated-support",
				},
			}
			_, err := icingaServer.UpdateUser(ctx, username, attrs)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		username := "test-user"
		t.Run("UserExists", func(t *testing.T) {
			err := icingaServer.DeleteUser(ctx, username)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("UserDoNotExist", func(t *testing.T) {
			err := icingaServer.DeleteUser(ctx, username)
			if err == nil || err.Error() != "No objects found." {
				t.Error(err)
			}
		})
	})
}
