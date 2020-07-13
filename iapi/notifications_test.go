package iapi

import (
	"testing"
)

func TestNotifications(t *testing.T) {
	icingaServer := Server{"root", ICINGA2_API_PASSWORD, "https://127.0.0.1:5665/v1", true, nil}
	// testHostName := "notification-test-host"
	// _, err := icingaServer.CreateHost(testHostName, "127.0.0.1", "hostalive", nil, nil, nil)
	// if err != nil {
	// 	t.Error(err)
	// }
	t.Run("Read", func(t *testing.T) {
		t.Run("TestGetValidNotification", func(t *testing.T) {
			name := "valid-notification"
			_, err := icingaServer.GetNotification(name)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("TestGetInvalidNotification", func(t *testing.T) {
			name := "invalid-notification"
			_, err := icingaServer.GetNotification(name)
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

			_, err := icingaServer.CreateNotification(notificationname, hostname, command, servicename, interval, nil, nil, nil)
			if err == nil {
				t.Error(err)
			}
		})

		t.Run("NotificationHostDoNotExist", func(t *testing.T) {
			hostname := "host-dne"
			notificationname := hostname + "!" + hostname
			command := "mail-host-notification"
			servicename := ""
			interval := 1800

			_, err := icingaServer.CreateNotification(notificationname, hostname, command, servicename, interval, nil, nil, nil)
			if err == nil {
				t.Error(err)
			}
		})

		t.Run("NotificationInvalidName", func(t *testing.T) {
			hostname := "host-dne"
			notificationname := "invalid-name"
			command := "mail-host-notification"
			servicename := ""
			interval := 1800

			_, err := icingaServer.CreateNotification(notificationname, hostname, command, servicename, interval, nil, nil, nil)
			if err == nil {
				t.Error(err)
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

			_, _ = icingaServer.CreateHost(hostname, "127.0.0.1", "hostalive", nil, nil, group)
			_, err := icingaServer.CreateNotification(notificationname, hostname, command, servicename, interval, users, nil, nil)

			if err == nil {
				t.Error(err)
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

			_, _ = icingaServer.CreateHost(hostname, "127.0.0.1", "hostalive", nil, nil, groups)
			_, _ = icingaServer.CreateUser(username, "user@example.com")
			_, err := icingaServer.CreateNotification(notificationname, hostname, command, servicename, interval, []string{username}, nil, nil)
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

			_, err := icingaServer.CreateNotification(notificationname, hostname, command, servicename, interval, []string{username}, nil, nil)
			if err == nil {
				t.Error(err)
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

			_, _ = icingaServer.CreateUser(username, "user@example.com")
			_, _ = icingaServer.CreateService(servicename, hostname, checkCommand, nil, nil)
			_, err := icingaServer.CreateNotification(notificationname, hostname, command, servicename, interval, []string{username}, nil, nil)

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

			_, _ = icingaServer.CreateUser(username, "user@example.com")
			_, _ = icingaServer.CreateService(servicename, hostname, checkCommand, nil, nil)
			_, err := icingaServer.CreateNotification(notificationName, hostname, command, servicename, interval, []string{username}, nil, nil)
			if err == nil {
				t.Error(err)
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("TestDeleteServiceNotification", func(t *testing.T) {
			hostname := "host"
			servicename := "test"
			notificationName := hostname + "!" + servicename + "!" + hostname + "-" + servicename
			err := icingaServer.DeleteNotification(notificationName)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("TestDeleteServiceNotificationDNE", func(t *testing.T) {

			hostname := "host"
			servicename := "test"
			notificationName := hostname + "!" + servicename + "!" + hostname + "-" + servicename

			err := icingaServer.DeleteNotification(notificationName)
			if err.Error() != "No objects found." {
				t.Error(err)
			}
		})

		t.Run("TestDeleteHostNotification", func(t *testing.T) {

			hostname := "host"
			notificationName := hostname + "!" + hostname
			username := "user"

			err := icingaServer.DeleteNotification(notificationName)
			_ = icingaServer.DeleteHost(hostname)
			_ = icingaServer.DeleteUser(username)

			if err != nil {
				t.Error(err)
			}
		})

		t.Run("TestDeleteHostNotificationDNE", func(t *testing.T) {

			hostname := "host"
			notificationName := hostname + "!" + hostname

			err := icingaServer.DeleteNotification(notificationName)

			if err.Error() != "No objects found." {
				t.Error(err)
			}
		})

		t.Run("TestDeleteNotificationNonAPI", func(t *testing.T) {

			notificationName := "mail-icingaadmin"

			err := icingaServer.DeleteNotification(notificationName)
			if err.Error() != "No objects found." {
				t.Error(err)
			}
		})
	})
>>>>>>> master
}
