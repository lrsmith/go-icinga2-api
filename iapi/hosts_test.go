package iapi

import "testing"

func TestGetValidHost(t *testing.T) {

	hostname := "c1-mysql-1"

	hosts, err := VagrantImage.GetHost(hostname)

	if err != nil {
		t.Errorf("Error : Failed to find %s : %s ) ", hostname, err)
	}

	if len(hosts) != 1 {
		t.Errorf("Error : Did not get expected number of results. Expected 1 got %d", len(hosts))
	}

	if hosts[0].Name != hostname {
		t.Errorf("Error : Did not get expected host. ( %s != %s )", hosts[0].Name, hostname)
	}

}

func TestGetInvalidHost(t *testing.T) {

	hostname := "c2-mysql-1"
	host, err := VagrantImage.GetHost(hostname)
	if err != nil && host != nil {
		t.Errorf("Error : Did not get empty list. ( %v : %s )", err, host)
	}

}

func TestCreateSimpleHost(t *testing.T) {

	hostname := "go-icinga2-api-1"
	IPAddress := "127.0.0.2"
	CheckCommand := "CheckItRealGood"
	_, err := VagrantImage.CreateHost(hostname, IPAddress, CheckCommand, nil)

	if err != nil {
		t.Errorf("Error : Failed to create host %s : %s", hostname, err)
	}

}

func TestCreateHostWithVariables(t *testing.T) {

	hostname := "go-icinga2-api-2"
	IPAddress := "127.0.0.3"
	CheckCommand := "CheckItRealGood"
	variables := make(map[string]string)

	variables["vars.os"] = "Linux"
	variables["vars.creator"] = "Terraform"

	_, err := VagrantImage.CreateHost(hostname, IPAddress, CheckCommand, variables)
	if err != nil {
		t.Errorf("Error : Failed to create host %s : %s", hostname, err)
	}

	// Delete host after creating it.
	deleteErr := VagrantImage.DeleteHost(hostname)
	if deleteErr != nil {
		t.Errorf("Error : Cleanup failed for host %s : %s", hostname, err)
	}
}

func TestDeleteHost(t *testing.T) {

	hostname := "go-icinga2-api-1"

	err := VagrantImage.DeleteHost(hostname)
	if err != nil {
		t.Errorf("Error : Failed to delete %s : %s", hostname, err)
	}

}

func TestDeleteHostDNE(t *testing.T) {
	hostname := "go-icinga2-api-1"
	err := VagrantImage.DeleteHost(hostname)
	if err.Error() != "No objects found." {
		t.Error(err)
	}
}
