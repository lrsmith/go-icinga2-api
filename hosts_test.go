package icinga2

import "testing"

func TestGetValidHost(t *testing.T) {

	hostname := "c1-mysql-1"

	host, err := VagrantImage.GetHost(hostname)

	if (err != nil) || (host == nil) {
		t.Errorf("Error : Failed to find %s : ( %s <> %s ) ", hostname, err, host)
	}

	if host[0].Name != hostname {
		t.Errorf("Error : Did not get expected hostname. ( %s != %s )", host[0].Name, hostname)
	}

}

func TestGetInvalidHost(t *testing.T) {

	hostname := "c2-mysql-1"
	host, err := VagrantImage.GetHost(hostname)
	if (err != nil) || (len(host) != 0) {
		t.Errorf("Error : Did not get empty list. ( %v : %s )", err, host)
	}

}

func TestCreateSimpleHost(t *testing.T) {

	hostname := "go-icinga2-api-1"
	IPAddress := "127.0.0.2"
	CheckCommand := "CheckItRealGood"
	err := VagrantImage.CreateHost(hostname, IPAddress, CheckCommand, nil)

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

	err := VagrantImage.CreateHost(hostname, IPAddress, CheckCommand, variables)
	if err != nil {
		t.Errorf("Error : Failed to create host %s : %s", hostname, err)
	}

	_ = VagrantImage.DeleteHost(hostname)

}

func TestDeleteHost(t *testing.T) {

	hostname := "go-icinga2-api-1"

	err := VagrantImage.DeleteHost(hostname)
	if err != nil {
		t.Errorf("Error : Failed to delete %s : %s", hostname, err)
	}

}

func TestDeleteNonExistentHost(t *testing.T) {
	hostname := "go-icinga2-api-1"
	err := VagrantImage.DeleteHost(hostname)
	if err != nil {
		t.Errorf("Error : Failed to delete %s : %s", hostname, err)
	}

}
