package iapi

import (
	"testing"
)

func TestHostgroups(t *testing.T) {
	icingaServer := Server{"root", ICINGA2_API_PASSWORD, "https://127.0.0.1:5665/v1", true, nil}
	t.Run("Create", func(t *testing.T) {
		t.Run("Hostgroup", func(t *testing.T) {
			name := "docker-servers"
			displayName := "Docker Host Servers"
			_, err := icingaServer.CreateHostgroup(name, displayName)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("Read", func(t *testing.T) {
		t.Run("ValidHostgroup", func(t *testing.T) {
			name := "linux-servers"
			_, err := icingaServer.GetHostgroup(name)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("InvalidHostgroup", func(t *testing.T) {
			name := "irix-servers"
			_, err := icingaServer.GetHostgroup(name)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("Update", func(t *testing.T) {
		hostgroupName := "someHostgroupName"
		firstDisplayName := "some Hostgroup Display Name"
		_, err := icingaServer.CreateHostgroup(hostgroupName, firstDisplayName)
		if err != nil {
			t.Error(err)
		}
		defer icingaServer.DeleteHostgroup(hostgroupName)

		secondDisplayName := "other hostgroup display name"
		params := &HostgroupParams{
			DisplayName: secondDisplayName,
		}
		updatedHostgroup, err := icingaServer.UpdateHostgroup(hostgroupName, params)
		if err != nil {
			t.Error(err)
		}
		if secondDisplayName != updatedHostgroup[0].Attrs.DisplayName {
			t.Errorf("expected display_name fields to not be equal, got equal")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		// Delete Hostgroup created via API. Should succeed
		t.Run("Hostgroup", func(t *testing.T) {
			name := "docker-servers"
			err := icingaServer.DeleteHostgroup(name)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("HostgroupNonAPI", func(t *testing.T) {
			name := "linux-servers"
			err := icingaServer.DeleteHostgroup(name)
			if err == nil {
				t.Error(err)
			}
		})

		t.Run("HostgroupDNE", func(t *testing.T) {
			name := "docker-servers"
			err := icingaServer.DeleteHostgroup(name)
			if err.Error() != "No objects found." {
				t.Error(err)
			}
		})
	})
}
