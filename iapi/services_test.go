package iapi

import (
	"context"
	"slices"
	"strings"
	"testing"
)

func TestServices(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}
	icingaServer, err := New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, ICINGA2_API_URL, ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	testHostName := "c1-mysql-1"
	_, err = icingaServer.CreateHost(ctx, testHostName, "127.0.0.1", "", "hostalive", nil, nil, nil, "")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		_ = icingaServer.DeleteHost(ctx, testHostName)
	}()

	t.Run("CreateService", func(t *testing.T) {
		// Try and create a service for a host that does not exist.
		// Should fail with an error about the host not existing.
		t.Run("HostDoNotExists", func(t *testing.T) {
			nonExistingHost := "c1-host-dne-1"
			servicename := "ssh"
			checkCommand := "ssh"

			_, err := icingaServer.CreateService(ctx, servicename, nonExistingHost, checkCommand, nil, nil)
			if err == nil {
				t.Error("ServiceHostDoNotExists: expected error returning, got nil")
			}
		})

		// Create a host and service via the API
		t.Run("HostAndService", func(t *testing.T) {
			servicename := "ssh"
			checkCommand := "ssh"

			_, err := icingaServer.CreateService(ctx, servicename, testHostName, checkCommand, nil, nil)
			if err != nil {
				t.Errorf("Error : Failed to create service %s!%s : %s", testHostName, servicename, err)
			}
		})

		t.Run("WithVariables", func(t *testing.T) {
			servicename := "nrpe"
			checkCommand := "nrpe"
			variables := make(map[string]string)
			variables["nrpe_command"] = "check_load"

			services, err := icingaServer.CreateService(ctx, servicename, testHostName, checkCommand, variables, nil)
			if err != nil {
				t.Errorf("Error : Failed to create service %s!%s : %s", testHostName, servicename, err)
			}

			vars := serviceVars(t, services, testHostName+"!"+servicename)
			if vars["nrpe_command"] != "check_load" {
				t.Errorf("expected variable nrpe_command to be check_load, got %v", vars["nrpe_command"])
			}
		})

		t.Run("WithTemplates", func(t *testing.T) {
			servicename := "nrpe-check"
			checkCommand := "nrpe"
			variables := make(map[string]string)
			variables["nrpe_command"] = "check_load"
			serviceTemplates := []string{"go-icinga2-api-test-service"}

			services, err := icingaServer.CreateService(ctx, servicename, testHostName, checkCommand, variables, serviceTemplates)
			if err != nil {
				t.Fatalf("Error : Failed to create service %s!%s : %s", testHostName, servicename, err)
			}

			name := testHostName + "!" + servicename
			for _, service := range services {
				if service.Name != name {
					continue
				}
				// Icinga resolves the imports of the imported templates too and
				// reports the service itself first
				expected := []string{servicename, "go-icinga2-api-test-service", "generic-service"}
				if !slices.Equal(service.Attrs.Templates, expected) {
					t.Errorf("expected templates %v, got %v", expected, service.Attrs.Templates)
				}
			}

			vars := serviceVars(t, services, name)
			if vars["nrpe_command"] != "check_load" {
				t.Errorf("expected variable nrpe_command to be check_load, got %v", vars["nrpe_command"])
			}
			// A variable of the template only reaches the service when the
			// template is imported instead of being stored as an attribute
			if vars["service_owner"] != "team" {
				t.Errorf("expected variable service_owner to be inherited from the template, got %v", vars["service_owner"])
			}
		})

		t.Run("WithUnknownTemplate", func(t *testing.T) {
			servicename := "unknown-template"
			checkCommand := "nrpe"
			serviceTemplates := []string{"go-icinga2-api-test-does-not-exist"}

			_, err := icingaServer.CreateService(ctx, servicename, testHostName, checkCommand, nil, serviceTemplates)
			if err == nil {
				t.Fatalf("expected the creation of %s to be rejected", servicename)
			}

			if !strings.Contains(err.Error(), "Import references unknown template") {
				t.Errorf("expected the reason of the rejection, got %s", err)
			}
		})

		// Test creating a host/service pair that already exists. Should get error about it already existing.
		t.Run("AlreadyExists", func(t *testing.T) {
			servicename := "ssh"
			checkCommand := "ssh"

			_, err = icingaServer.CreateService(ctx, servicename, testHostName, checkCommand, nil, nil)
			if err == nil {
				t.Error("TestCreateServiceAlreadyExists: expected error returning, got nil")
			}
		})
	})

	t.Run("ReadService", func(t *testing.T) {
		t.Run("ValidService", func(t *testing.T) {
			servicename := "ssh"
			_, err := icingaServer.GetService(ctx, servicename, testHostName)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("InvalidService", func(t *testing.T) {
			servicename := "foo"
			_, err := icingaServer.GetService(ctx, servicename, testHostName)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("UpdateService", func(t *testing.T) {
		t.Run("ValidService", func(t *testing.T) {
			servicename := "nrpe"
			attrs := ServiceAttrs{
				CheckCommand: "nrpe",
				Vars: map[string]string{
					"vars.nrpe_command": "check_updated",
				},
			}
			_, err := icingaServer.UpdateService(ctx, servicename, testHostName, attrs)
			if err != nil {
				t.Error(err)
			}
		})

		// The templates of a service are immutable, Icinga reports a per object
		// error while answering 200
		t.Run("TemplatesAreRejected", func(t *testing.T) {
			servicename := "nrpe-check"
			attrs := ServiceAttrs{
				CheckCommand: "nrpe",
				Templates:    []string{"go-icinga2-api-test-service"},
			}
			_, err := icingaServer.UpdateService(ctx, servicename, testHostName, attrs)
			if err == nil {
				t.Fatal("expected the update to be rejected")
			}
			if !strings.Contains(err.Error(), "Attribute cannot be modified") {
				t.Errorf("expected the reason of the rejection, got %s", err)
			}
		})
	})

	t.Run("DeleteService", func(t *testing.T) {
		// Delete a service which was create via the API.
		// Should not get an error
		t.Run("HostAndService", func(t *testing.T) {
			servicename := "ssh"

			err := icingaServer.DeleteService(ctx, servicename, testHostName)
			if err != nil {
				t.Error(err)
			}
		})

		// Try and delete a service, where the host does not exists.
		// Should get an error abot no object found
		t.Run("ServiceHostDoNotExists", func(t *testing.T) {
			hostname := "c1-test-1"
			servicename := "ssh"

			err := icingaServer.DeleteService(ctx, servicename, hostname)
			if err == nil || err.Error() != "No objects found." {
				t.Error(err)
			}
		})

		// Try and delete a service, where the host exists but the service does not.
		// Should get an error abot no object found
		t.Run("ServiceDoNotExists", func(t *testing.T) {
			servicename := "foo"
			err := icingaServer.DeleteService(ctx, servicename, testHostName)
			if err == nil || err.Error() != "No objects found." {
				t.Error(err)
			}
		})

		// Services that were not created via the API, cannot be deleted via the API
		// Should get an error about not being created via the API
		t.Run("ServiceNonAPI", func(t *testing.T) {
			hostname := "docker-icinga2"
			servicename := "random-001"

			err := icingaServer.DeleteService(ctx, servicename, hostname)
			if err == nil || err.Error() != "No objects found." {
				t.Error(err)
			}
		})
	})

}

// serviceVars returns the resolved variables of a service.
func serviceVars(t *testing.T, services []ServiceStruct, name string) map[string]interface{} {
	t.Helper()

	for _, service := range services {
		if service.Name != name {
			continue
		}
		vars, ok := service.Attrs.Vars.(map[string]interface{})
		if !ok {
			t.Fatalf("expected variables on %s, got %v", name, service.Attrs.Vars)
		}
		return vars
	}

	t.Fatalf("service %s not found", name)
	return nil
}
