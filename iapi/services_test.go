package iapi

import (
	"testing"
)

func TestServices(t *testing.T) {
	icingaServer := Server{"root", ICINGA2_API_PASSWORD, "https://127.0.0.1:5665/v1", true, nil}
	testHostName := "c1-mysql-1"
	_, err := icingaServer.CreateHost(testHostName, "127.0.0.1", "hostalive", nil, nil, nil)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		_ = icingaServer.DeleteHost(testHostName)
	}()

	t.Run("CreateService", func(t *testing.T) {
		// Try and create a service for a host that does not exist.
		// Should fail with an error about the host not existing.
		t.Run("HostDoNotExists", func(t *testing.T) {
			nonExistingHost := "c1-host-dne-1"
			servicename := "ssh"
			checkCommand := "ssh"

			_, err := icingaServer.CreateService(servicename, nonExistingHost, checkCommand, nil, nil)
			if err == nil {
				t.Error("ServiceHostDoNotExists: expected error returning, got nil")
			}
		})

		// Create a host and service via the API
		t.Run("HostAndService", func(t *testing.T) {
			servicename := "ssh"
			checkCommand := "ssh"

			_, err := icingaServer.CreateService(servicename, testHostName, checkCommand, nil, nil)
			if err != nil {
				t.Errorf("Error : Failed to create service %s!%s : %s", testHostName, servicename, err)
			}
		})

		t.Run("WithVariables", func(t *testing.T) {
			servicename := "nrpe"
			checkCommand := "nrpe"
			variables := make(map[string]string)
			variables["vars.nrpe_command"] = "check_load"

			_, err := icingaServer.CreateService(servicename, testHostName, checkCommand, variables, nil)
			if err != nil {
				t.Errorf("Error : Failed to create service %s!%s : %s", testHostName, servicename, err)
			}
		})

		t.Run("WithTemplates", func(t *testing.T) {
			servicename := "nrpe-check"
			checkCommand := "nrpe"
			variables := make(map[string]string)
			variables["vars.nrpe_command"] = "check_load"
			serviceTemplates := []string{"generic-service", "holiwi"}

			_, err = icingaServer.CreateService(servicename, testHostName, checkCommand, variables, serviceTemplates)
			if err != nil {
				t.Errorf("Error : Failed to create service %s!%s : %s", testHostName, servicename, err)
			}
		})

		// Test creating a host/service pair that already exists. Should get error about it already existing.
		t.Run("AlreadyExists", func(t *testing.T) {
			servicename := "ssh"
			checkCommand := "ssh"

			_, err = icingaServer.CreateService(servicename, testHostName, checkCommand, nil, nil)
			if err == nil {
				t.Error("TestCreateServiceAlreadyExists: expected error returning, got nil")
			}
		})
	})

	t.Run("ReadService", func(t *testing.T) {
		t.Run("ValidService", func(t *testing.T) {
			servicename := "ssh"
			_, err := icingaServer.GetService(servicename, testHostName)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("InvalidService", func(t *testing.T) {
			servicename := "foo"
			_, err := icingaServer.GetService(servicename, testHostName)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("DeleteService", func(t *testing.T) {
		// Delete a service which was create via the API.
		// Should not get an error
		t.Run("HostAndService", func(t *testing.T) {
			servicename := "ssh"

			err := icingaServer.DeleteService(servicename, testHostName)
			if err != nil {
				t.Error(err)
			}
		})

		// Try and delete a service, where the host does not exists.
		// Should get an error abot no object found
		t.Run("ServiceHostDoNotExists", func(t *testing.T) {
			hostname := "c1-test-1"
			servicename := "ssh"

			err := icingaServer.DeleteService(servicename, hostname)
			if err.Error() != "No objects found." {
				t.Error(err)
			}
		})

		// Try and delete a service, where the host exists but the service does not.
		// Should get an error abot no object found
		t.Run("ServiceDoNotExists", func(t *testing.T) {
			servicename := "foo"
			err := icingaServer.DeleteService(servicename, testHostName)
			if err.Error() != "No objects found." {
				t.Error(err)
			}
		})

		// Services that were not created via the API, cannot be deleted via the API
		// Should get an error about not being created via the API
		t.Run("ServiceNonAPI", func(t *testing.T) {
			hostname := "docker-icinga2"
			servicename := "random-001"

			err := icingaServer.DeleteService(servicename, hostname)
			if err.Error() != "No objects found." {
				t.Error(err)
			}
		})
	})

}
