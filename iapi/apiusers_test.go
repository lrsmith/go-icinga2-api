package iapi

import (
	"slices"
	"testing"
)

func TestApiUsers(t *testing.T) {
	icingaServer := Server{ICINGA2_API_USER, ICINGA2_API_PASSWORD, ICINGA2_API_URL, ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0, nil}
	t.Run("Create", func(t *testing.T) {
		t.Run("ApiUser", func(t *testing.T) {
			name := "director"
			password := "director-password"
			permissions := []string{"*"}
			_, err := icingaServer.CreateApiUser(name, password, "", permissions)
			if err != nil {
				t.Error(err)
			}
		})
		t.Run("ApiUserWithCNs", func(t *testing.T) {
			name := "director-client-cn"
			password := ""
			clientCN := "director"
			permissions := []string{"*"}
			_, err := icingaServer.CreateApiUser(name, password, clientCN, permissions)
			if err != nil {
				t.Error(err)
			}
		})

	})

	t.Run("Read", func(t *testing.T) {
		t.Run("ValidApiUser", func(t *testing.T) {
			name := "director"
			_, err := icingaServer.GetApiUser(name)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("InvalidApiUser", func(t *testing.T) {
			name := "some-nonexistent-user"
			_, err := icingaServer.GetApiUser(name)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("Update", func(t *testing.T) {
		apiUserName := "someApiUserName"
		password := "secret"
		permissions := []string{"*"}
		_, err := icingaServer.CreateApiUser(apiUserName, password, "", permissions)
		if err != nil {
			t.Error(err)
		}
		defer icingaServer.DeleteApiUser(apiUserName)

		secondPermissions := []string{"objects/query/Host", "objects/query/Service"}
		params := &ApiUserAttrs{
			Permissions: secondPermissions,
		}
		updatedApiUser, err := icingaServer.UpdateApiUser(apiUserName, params)
		if err != nil {
			t.Error(err)
		}
		if !slices.Equal(secondPermissions, updatedApiUser[0].Attrs.Permissions) {
			t.Errorf("expected permissions are not found to be equal.")
		}
	})

	t.Run("Exists", func(t *testing.T) {
		t.Run("ApiUserFound", func(t *testing.T) {
			name := "director"
			exists, err := icingaServer.ApiUserExists(name)
			if err != nil {
				t.Error(err)
			}

			if !exists {
				t.Error("apiuser must exist")
			}
		})

		t.Run("ApiUserNotFound", func(t *testing.T) {
			name := "some-nonexistent-user"
			exists, err := icingaServer.ApiUserExists(name)
			if err != nil {
				t.Error(err)
			}
			if exists {
				t.Error("apiuser must not exist")
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		// Delete ApiUser created via API. Should succeed
		t.Run("ApiUser", func(t *testing.T) {
			name := "director"
			err := icingaServer.DeleteApiUser(name)
			if err != nil {
				t.Error(err)
			}
		})
		t.Run("ApiUserWithCNs", func(t *testing.T) {
			name := "director-client-cn"
			err := icingaServer.DeleteApiUser(name)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("ApiUserNonAPI", func(t *testing.T) {
			name := "root" // Assuming this wasn't created via API or doesn't exist
			err := icingaServer.DeleteApiUser(name)
			if err == nil {
				t.Error(err)
			}
		})

		t.Run("ApiUserDNE", func(t *testing.T) {
			name := "director" // Already deleted
			err := icingaServer.DeleteApiUser(name)
			if err.Error() != "No objects found." {
				t.Error(err)
			}
		})
	})
}
