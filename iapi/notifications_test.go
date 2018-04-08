package iapi

import (
	"strings"
	"testing"
)

func TestGetValidNotification(t *testing.T) {

	name := "valid-notification"
	_, err := Icinga2_Server.GetNotification(name)

	if err != nil {
		t.Error(err)
	}
}

func TestGetInvalidNotification(t *testing.T) {

	name := "invalid-notification"

	_, err := Icinga2_Server.GetNotification(name)

	if err != nil {
		t.Error(err)
	}
}

func TestCreateNotificationCommandDNE(t *testing.T) {

	hostname := "host"
	notificationname := hostname+"!"+hostname
	command := "invalid-command"
	servicename := ""
	interval := 1800

	_, err := Icinga2_Server.CreateNotification(notificationname, hostname, command, servicename, interval, nil, nil, nil)

	if !strings.Contains(err.Error(), "type 'NotificationCommand' does not exist.") {
		t.Error(err)
	}
}

// func TestCreateNotificationHostDNE
// Try and create a notification for a host that does not exist.
// Should fail with an error about the host not existing.
func TestCreateNotificationHostDNE(t *testing.T) {

	hostname := "host-dne"
	notificationname := hostname+"!"+hostname
	command := "mail-host-notification"
	servicename := ""
	interval := 1800

	_, err := Icinga2_Server.CreateNotification(notificationname, hostname, command, servicename, interval, nil, nil, nil)

	if !strings.Contains(err.Error(), "type 'Host' does not exist.") {
		t.Error(err)
	}
}

func TestCreateNotificationInvalidName(t *testing.T) {

	hostname := "host-dne"
	notificationname := "invalid-name"
	command := "mail-host-notification"
	servicename := ""
	interval := 1800

	_, err := Icinga2_Server.CreateNotification(notificationname, hostname, command, servicename, interval, nil, nil, nil)

	if !strings.Contains(err.Error(), "Invalid Notification name.") {
		t.Error(err)
	}
}

func TestCreateNotificationUserDNE(t *testing.T) {

	hostname := "host"
	group := []string{"linux-servers"}
	notificationname := hostname+"!"+hostname
	command := "mail-host-notification"
	servicename := ""
	interval := 1800
	users := []string{"user-dne"}


    _, _ = Icinga2_Server.CreateHost(hostname, "127.0.0.1", "hostalive", nil, nil, group)
	_, err := Icinga2_Server.CreateNotification(notificationname, hostname, command, servicename, interval, users, nil, nil)

	if !strings.Contains(err.Error(), "type 'User' does not exist.") {
		t.Error(err)
	}
}

func TestCreateHostNotification(t *testing.T) {

	hostname := "host"
	groups := []string{"linux-servers"}
	notificationname := hostname+"!"+hostname
	command := "mail-host-notification"
	servicename := ""
	interval := 1800
	username := "user"

    _, _ = Icinga2_Server.CreateHost(hostname, "127.0.0.1", "hostalive", nil, nil, groups)
    _, _ = Icinga2_Server.CreateUser(username, "user@example.com")
	_, err := Icinga2_Server.CreateNotification(notificationname, hostname, command, servicename, interval, []string{username}, nil, nil)

	if err != nil {
		t.Error(err)
	}
}

func TestCreateHostNotificationAlreadyExists(t *testing.T) {

	hostname := "host"
	notificationname := hostname+"!"+hostname
	command := "mail-host-notification"
	servicename := ""
	interval := 1800
	username := "user"

	_, err := Icinga2_Server.CreateNotification(notificationname, hostname, command, servicename, interval, []string{username}, nil, nil)

	if !strings.HasSuffix(err.Error(), " already exists.") {
		t.Error(err)
	}
}

func TestCreateServiceNotification(t *testing.T) {

	hostname := "host"
	servicename := "test"
	check_command := "random"
	notificationname := hostname+"!"+servicename+"!"+hostname+"-"+servicename
	command := "mail-service-notification"
	interval := 1800
	username := "user"

    _, _ = Icinga2_Server.CreateUser(username, "user@example.com")
    _, _ = Icinga2_Server.CreateService(servicename, hostname, check_command,nil)
	_, err := Icinga2_Server.CreateNotification(notificationname, hostname, command, servicename, interval, []string{username}, nil, nil)

	if err != nil {
		t.Error(err)
	}
}

func TestCreateServiceNotificationAlreadyExists(t *testing.T) {

	hostname := "host"
	servicename := "test"
	check_command := "random"
	notificationname := hostname+"!"+servicename+"!"+hostname+"-"+servicename
	command := "mail-service-notification"
	interval := 1800
	username := "user"

    _, _ = Icinga2_Server.CreateUser(username, "user@example.com")
    _, _ = Icinga2_Server.CreateService(servicename, hostname, check_command, nil)
	_, err := Icinga2_Server.CreateNotification(notificationname, hostname, command, servicename, interval, []string{username}, nil, nil)

	if !strings.HasSuffix(err.Error(), " already exists.") {
		t.Error(err)
	}
}

func TestDeleteServiceNotification(t *testing.T) {

	hostname := "host"
	servicename := "test"
	notificationname := hostname+"!"+servicename+"!"+hostname+"-"+servicename

	err := Icinga2_Server.DeleteNotification(notificationname)

	if err != nil {
		t.Error(err)
	}
}

func TestDeleteServiceNotificationDNE(t *testing.T) {

	hostname := "host"
	servicename := "test"
	notificationname := hostname+"!"+servicename+"!"+hostname+"-"+servicename

	err := Icinga2_Server.DeleteNotification(notificationname)

	if err.Error() != "No objects found." {
		t.Error(err)
	}
}

func TestDeleteHostNotification(t *testing.T) {

	hostname := "host"
	notificationname := hostname+"!"+hostname
	username := "user"

	err := Icinga2_Server.DeleteNotification(notificationname)
    _ = Icinga2_Server.DeleteHost(hostname)
    _ = Icinga2_Server.DeleteUser(username)

	if err != nil {
		t.Error(err)
	}
}

func TestDeleteHostNotificationDNE(t *testing.T) {

	hostname := "host"
	notificationname := hostname+"!"+hostname

	err := Icinga2_Server.DeleteNotification(notificationname)

	if err.Error() != "No objects found." {
		t.Error(err)
	}
}

// func TestDeleteNotificationNonAPI
// Notifications that were not created via the API, cannot be deleted via the API
// Should get an error about not being created via the API
func TestDeleteNotificationNonAPI(t *testing.T) {

	notificationname := "mail-icingaadmin"

	err := Icinga2_Server.DeleteNotification(notificationname)
	if err.Error() != "No objects found." {
		t.Error(err)
	}
}
