package iapi

import (
	"context"
	"slices"
	"strings"
	"testing"
)

func TestNotifications(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}
	icingaServer, err := New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, ICINGA2_API_URL, ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	t.Run("Read", func(t *testing.T) {
		t.Run("TestGetValidNotification", func(t *testing.T) {
			name := "valid-notification"
			_, err := icingaServer.GetNotification(ctx, name)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("TestGetInvalidNotification", func(t *testing.T) {
			name := "invalid-notification"
			_, err := icingaServer.GetNotification(ctx, name)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("NotificationCommandDoNotExist", func(t *testing.T) {
			hostname := "host"
			notificationname := hostname + "!" + hostname
			command := "invalid-command"
			servicename := ""
			interval := 1800

			_, err := icingaServer.CreateNotification(ctx, notificationname, hostname, command, servicename, interval, nil, nil, nil)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})

		t.Run("NotificationHostDoNotExist", func(t *testing.T) {
			hostname := "host-dne"
			notificationname := hostname + "!" + hostname
			command := "mail-host-notification"
			servicename := ""
			interval := 1800

			_, err := icingaServer.CreateNotification(ctx, notificationname, hostname, command, servicename, interval, nil, nil, nil)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})

		t.Run("NotificationInvalidName", func(t *testing.T) {
			hostname := "host-dne"
			notificationname := "invalid-name"
			command := "mail-host-notification"
			servicename := ""
			interval := 1800

			_, err := icingaServer.CreateNotification(ctx, notificationname, hostname, command, servicename, interval, nil, nil, nil)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})

		t.Run("NotificationUserDoNotExist", func(t *testing.T) {
			hostname := "host"
			group := []string{"linux-servers"}
			notificationname := hostname + "!" + hostname
			command := "mail-host-notification"
			servicename := ""
			interval := 1800
			users := []string{"user-dne"}

			_, _ = icingaServer.CreateHost(ctx, hostname, "127.0.0.1", "", "hostalive", nil, nil, group, "")
			_, err := icingaServer.CreateNotification(ctx, notificationname, hostname, command, servicename, interval, users, nil, nil)

			if err == nil {
				t.Error("expected error, got nil")
			}
		})

		t.Run("HostNotification", func(t *testing.T) {
			hostname := "host"
			groups := []string{"linux-servers"}
			notificationname := hostname + "!" + hostname
			command := "mail-host-notification"
			servicename := ""
			interval := 1800
			username := "user"

			_, _ = icingaServer.CreateHost(ctx, hostname, "127.0.0.1", "", "hostalive", nil, nil, groups, "")
			_, _ = icingaServer.CreateUser(ctx, username, "user@example.com", nil)
			_, err := icingaServer.CreateNotification(ctx, notificationname, hostname, command, servicename, interval, []string{username}, nil, nil)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("HostNotificationAlreadyExists", func(t *testing.T) {
			hostname := "host"
			notificationname := hostname + "!" + hostname
			command := "mail-host-notification"
			servicename := ""
			interval := 1800
			username := "user"

			_, err := icingaServer.CreateNotification(ctx, notificationname, hostname, command, servicename, interval, []string{username}, nil, nil)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})

		t.Run("ServiceNotification", func(t *testing.T) {
			hostname := "host"
			servicename := "test"
			checkCommand := "random"
			notificationname := hostname + "!" + servicename + "!" + hostname + "-" + servicename
			command := "mail-service-notification"
			interval := 1800
			username := "user"

			_, _ = icingaServer.CreateUser(ctx, username, "user@example.com", nil)
			_, _ = icingaServer.CreateService(ctx, servicename, hostname, checkCommand, nil, nil)
			_, err := icingaServer.CreateNotification(ctx, notificationname, hostname, command, servicename, interval, []string{username}, nil, nil)

			if err != nil {
				t.Error(err)
			}
		})

		t.Run("ServiceNotificationAlreadyExists", func(t *testing.T) {
			hostname := "host"
			servicename := "test"
			checkCommand := "random"
			notificationName := hostname + "!" + servicename + "!" + hostname + "-" + servicename
			command := "mail-service-notification"
			interval := 1800
			username := "user"

			_, _ = icingaServer.CreateUser(ctx, username, "user@example.com", nil)
			_, _ = icingaServer.CreateService(ctx, servicename, hostname, checkCommand, nil, nil)
			_, err := icingaServer.CreateNotification(ctx, notificationName, hostname, command, servicename, interval, []string{username}, nil, nil)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})

		t.Run("WithTemplates", func(t *testing.T) {
			hostname := "host"
			notificationname := hostname + "!" + hostname + "-templated"
			username := "user"
			templates := []string{"go-icinga2-api-test-notification"}

			_, _ = icingaServer.CreateUser(ctx, username, "user@example.com", nil)

			// The command and the interval are left to the templates, which only
			// works when they are imported instead of being stored as attributes
			notifications, err := icingaServer.CreateNotification(ctx, notificationname, hostname, "", "", 0, []string{username}, nil, templates)
			if err != nil {
				t.Fatal(err)
			}

			for _, notification := range notifications {
				if notification.Name != notificationname {
					continue
				}
				if !slices.Contains(notification.Attrs.Templates, templates[0]) {
					t.Errorf("expected %s to import %s, got %v", notificationname, templates[0], notification.Attrs.Templates)
				}
				if notification.Attrs.Command != "mail-host-notification" {
					t.Errorf("expected the command to be inherited from the template, got %q", notification.Attrs.Command)
				}
				if notification.Attrs.Interval != 900 {
					t.Errorf("expected the interval to be inherited from the template, got %d", notification.Attrs.Interval)
				}
			}

			err = icingaServer.DeleteNotification(ctx, notificationname)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("WithUnknownTemplate", func(t *testing.T) {
			hostname := "host"
			notificationname := hostname + "!" + hostname + "-unknown-template"
			templates := []string{"go-icinga2-api-test-does-not-exist"}

			_, err := icingaServer.CreateNotification(ctx, notificationname, hostname, "mail-host-notification", "", 1800, nil, nil, templates)
			if err == nil {
				t.Fatalf("expected the creation of %s to be rejected", notificationname)
			}

			if !strings.Contains(err.Error(), "Import references unknown template") {
				t.Errorf("expected the reason of the rejection, got %s", err)
			}
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("HostNotification", func(t *testing.T) {
			hostname := "host"
			notificationname := hostname + "!" + hostname
			attrs := NotificationAttrs{
				Command:  "mail-host-notification",
				Interval: 3600,
				Vars: map[string]string{
					"vars.custom_field": "updated",
				},
			}
			notifications, err := icingaServer.UpdateNotification(ctx, notificationname, attrs)
			if err != nil {
				t.Fatal(err)
			}

			for _, notification := range notifications {
				if notification.Name != notificationname {
					continue
				}
				if notification.Attrs.Interval != 3600 {
					t.Errorf("expected the interval to be updated to 3600, got %d", notification.Attrs.Interval)
				}
			}
		})

		// The service name and the templates of a notification are immutable,
		// sending them along makes Icinga reject the whole update
		t.Run("HostNotificationWithImmutableAttrs", func(t *testing.T) {
			hostname := "host"
			notificationname := hostname + "!" + hostname
			attrs := NotificationAttrs{
				Command:   "mail-host-notification",
				Interval:  7200,
				Templates: []string{"mail-host-notification"},
			}
			notifications, err := icingaServer.UpdateNotification(ctx, notificationname, attrs)
			if err != nil {
				t.Fatal(err)
			}

			for _, notification := range notifications {
				if notification.Name != notificationname {
					continue
				}
				if notification.Attrs.Interval != 7200 {
					t.Errorf("expected the interval to be updated to 7200, got %d", notification.Attrs.Interval)
				}
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("TestDeleteServiceNotification", func(t *testing.T) {
			hostname := "host"
			servicename := "test"
			notificationName := hostname + "!" + servicename + "!" + hostname + "-" + servicename
			err := icingaServer.DeleteNotification(ctx, notificationName)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("TestDeleteServiceNotificationDNE", func(t *testing.T) {
			hostname := "host"
			servicename := "test"
			notificationName := hostname + "!" + servicename + "!" + hostname + "-" + servicename

			err := icingaServer.DeleteNotification(ctx, notificationName)
			if err == nil || err.Error() != "No objects found." {
				t.Error(err)
			}
		})

		t.Run("TestDeleteHostNotification", func(t *testing.T) {
			hostname := "host"
			notificationName := hostname + "!" + hostname
			username := "user"

			err := icingaServer.DeleteNotification(ctx, notificationName)
			_ = icingaServer.DeleteHost(ctx, hostname)
			_ = icingaServer.DeleteUser(ctx, username)

			if err != nil {
				t.Error(err)
			}
		})

		t.Run("TestDeleteHostNotificationDNE", func(t *testing.T) {
			hostname := "host"
			notificationName := hostname + "!" + hostname

			err := icingaServer.DeleteNotification(ctx, notificationName)

			if err == nil || err.Error() != "No objects found." {
				t.Error(err)
			}
		})

		t.Run("TestDeleteNotificationNonAPI", func(t *testing.T) {
			notificationName := "mail-icingaadmin"

			err := icingaServer.DeleteNotification(ctx, notificationName)
			if err == nil || err.Error() != "No objects found." {
				t.Error(err)
			}
		})
	})
}
