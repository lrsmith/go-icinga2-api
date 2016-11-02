package iapi

import "testing"

var VagrantImage = Server{"root", "icinga", "https://192.168.33.5:5665/v1", true, nil}
var VagrantImageBadPassword = Server{"root", "icinga2", "https://192.168.33.5:5665/v1", true, nil}

func TestConnect(t *testing.T) {

	VagrantImage.Connect()

	if VagrantImage.httpClient == nil {
		t.Errorf("Failed to succesfull connect to Icinga Server")
	}
}

func TestConnectWithBadCredential(t *testing.T) {

	VagrantImageBadPassword.Connect()
	if VagrantImageBadPassword.httpClient != nil {
		t.Errorf("Did not fail with bad credentials")
	}
}

func TestNewAPIRequest(t *testing.T) {

	result, _ := VagrantImage.NewAPIRequest("GET", "/status", nil)

	if result.Status != "200 OK" {
		t.Errorf("%s", result.Status)
	}

	//fmt.Printf("%v", result.Results)

}
