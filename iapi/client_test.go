package iapi

import (
	"os"
	"testing"
)

var ICINGA2_API_PASSWORD = os.Getenv("ICINGA2_API_PASSWORD")

var Icinga2_Server = Server{"root", ICINGA2_API_PASSWORD, "https://127.0.0.1:5665/v1", true, nil}
//var Icinga2_Server = Server{"icinga-test", "icinga", "https://127.0.0.1:5665/v1", true, nil}

func TestConnect(t *testing.T) {

	var Icinga2_Server = Server{"icinga-test", "icinga", "https://127.0.0.1:5665/v1", true, nil}
	Icinga2_Server.Connect()

	if Icinga2_Server.httpClient == nil {
		t.Errorf("Failed to succesfully connect to Icinga Server")
	}
}

func TestConnectServerUnavailable(t *testing.T) {

	var Icinga2_Server = Server{"icinga-test", "icinga", "https://127.0.0.1:4665/v1", true, nil}
	err := Icinga2_Server.Connect()

	if err == nil {
		t.Errorf("Error : Did not get error connecting to unavailable server.")
	}
}

func TestConnectWithBadCredential(t *testing.T) {

	var Icinga2_Server = Server{"unknownUser", "unknownPW", "https://127.0.0.1:5665/v1", true, nil}
	err := Icinga2_Server.Connect()
	if err != nil {
		t.Errorf("Did not fail with bad credentials : %s", err)
	}
}

func TestNewAPIRequest(t *testing.T) {

	result, _ := Icinga2_Server.NewAPIRequest("GET", "/status", nil)

	if result.Code != 200 {
		t.Errorf("%s", result.Status)
	}
}

func TestConnectServerBadURINoVersion(t *testing.T) {

	var Icinga2_Server = Server{"root", ICINGA2_API_PASSWORD, "https://127.0.0.1:5665", true, nil}
	result, _ := Icinga2_Server.NewAPIRequest("GET", "/status", nil)

	if result.Code != 404 {
		t.Errorf("Error : Did not get expected 404 error connection to bad URI, with no version.")
	}
}
