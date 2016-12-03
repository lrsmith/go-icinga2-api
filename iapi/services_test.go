package iapi

import "testing"

func TestGetValidService(t *testing.T) {

	hostname := "c1-mysql-1"
	servicename := "ssh"

	services, err := VagrantImage.GetService(servicename, hostname)

	if err != nil {
		t.Errorf("Error : Failed to find %s for %s : ( %s <> %v ) ", servicename, hostname, err, services)
	}

	if len(services) != 1 {
		t.Errorf("Error : Did not get expected number of results. Expected 1 got %d", len(services))
	}

	if services[0].Name != hostname+"!"+servicename {
		t.Errorf("Error : Did not get expected service. ( %s != %s!%s )", services[0].Name, hostname, servicename)
	}

}

func TestGetInvalidService(t *testing.T) {

	hostname := "c1-mysql-1"
	servicename := "foo"

	services, err := VagrantImage.GetService(servicename, hostname)

	if err != nil {
		t.Errorf("Error : Failed to find %s for %s : ( %s <> %v ) ", servicename, hostname, err, services)
	}

	if len(services) != 0 {
		t.Errorf("Error : Did not get expected number of results. Expected 0 got %d", len(services))
	}

}

func TestCreateService(t *testing.T) {

	hostname := "c1-test-1"
	servicename := "ssh"
	check_command := "ssh"

	_, _ = VagrantImage.CreateHost(hostname, "127.0.0.1", "hostalive", nil)

	_, err := VagrantImage.CreateService(servicename, hostname, check_command)

	if err != nil {
		t.Errorf("Error : Failed to create service %s!%s : %s", hostname, servicename, err)
	}

}

func TestDeleteService(t *testing.T) {

	hostname := "c1-test-1"
	servicename := "ssh"

	err := VagrantImage.DeleteService(servicename, hostname)
	if err != nil {
		_ = VagrantImage.DeleteHost(hostname)
		t.Errorf("Error : Failed to delete %s!%s : %s", hostname, servicename, err)
	}

	_ = VagrantImage.DeleteHost(hostname)

}

/*
func TestDeleteNonExistentHost(t *testing.T) {
	hostname := "go-icinga2-api-1"
	err := VagrantImage.DeleteHost(hostname)
	if err != nil {
		t.Errorf("Error : Failed to delete %s : %s", hostname, err)
	}

}
*/
