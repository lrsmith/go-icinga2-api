package iapi

import (
	"context"
	"testing"
)

func TestHostgroups(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}
	icingaServer, err := New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, ICINGA2_API_URL, ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		t.Run("Hostgroup", func(t *testing.T) {
			name := "docker-servers"
			displayName := "Docker Host Servers"
			_, err := icingaServer.CreateHostgroup(ctx, name, displayName, "")
			if err != nil {
				t.Error(err)
			}
		})
		t.Run("HostgroupWithZone", func(t *testing.T) {
			name := "docker-servers-zone"
			displayName := "Docker Host Servers on zone"
			_, err := icingaServer.CreateHostgroup(ctx, name, displayName, "master")
			if err != nil {
				t.Error(err)
			}
		})

	})

	t.Run("Read", func(t *testing.T) {
		t.Run("ValidHostgroup", func(t *testing.T) {
			name := "linux-servers"
			_, err := icingaServer.GetHostgroup(ctx, name)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("InvalidHostgroup", func(t *testing.T) {
			name := "irix-servers"
			_, err := icingaServer.GetHostgroup(ctx, name)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("Update", func(t *testing.T) {
		hostgroupName := "someHostgroupName"
		firstDisplayName := "some Hostgroup Display Name"
		_, err := icingaServer.CreateHostgroup(ctx, hostgroupName, firstDisplayName, "")
		if err != nil {
			t.Fatalf("failed to create hostgroup: %v", err)
		}
		defer icingaServer.DeleteHostgroup(ctx, hostgroupName)

		secondDisplayName := "other hostgroup display name"
		attrs := HostgroupAttrs{
			DisplayName: secondDisplayName,
		}
		updatedHostgroup, err := icingaServer.UpdateHostgroup(ctx, hostgroupName, attrs)
		if err != nil {
			t.Fatalf("failed to update hostgroup: %v", err)
		}
		if len(updatedHostgroup) == 0 {
			t.Fatalf("expected non-empty updated hostgroup slice")
		}
		if secondDisplayName != updatedHostgroup[0].Attrs.DisplayName {
			t.Errorf("expected display_name fields to be equal, got different")
		}
	})

	t.Run("Exists", func(t *testing.T) {
		t.Run("HostgroupFound", func(t *testing.T) {
			name := "docker-servers"
			exists, err := icingaServer.HostgroupExists(ctx, name)
			if err != nil {
				t.Error(err)
			}

			if !exists {
				t.Error("host group must exist")
			}
		})

		t.Run("HostgroupNotFound", func(t *testing.T) {
			name := "irix-servers"
			exists, err := icingaServer.HostgroupExists(ctx, name)
			if err != nil {
				t.Error(err)
			}
			if exists {
				t.Error("host group must not exist")
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		// Delete Hostgroup created via API. Should succeed
		t.Run("Hostgroup", func(t *testing.T) {
			name := "docker-servers"
			err := icingaServer.DeleteHostgroup(ctx, name)
			if err != nil {
				t.Error(err)
			}
		})
		t.Run("HostgroupWithZone", func(t *testing.T) {
			name := "docker-servers-zone"
			err := icingaServer.DeleteHostgroup(ctx, name)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("HostgroupNonAPI", func(t *testing.T) {
			name := "linux-servers"
			err := icingaServer.DeleteHostgroup(ctx, name)
			if err == nil {
				t.Error("expected error deleting non-api hostgroup, got nil")
			}
		})

		t.Run("HostgroupDNE", func(t *testing.T) {
			name := "docker-servers"
			err := icingaServer.DeleteHostgroup(ctx, name)
			if err == nil || err.Error() != "No objects found." {
				t.Error(err)
			}
		})
	})
}
