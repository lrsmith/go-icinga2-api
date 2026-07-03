package iapi

import (
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/jarcoal/httpmock"
)

var ICINGA2_API_USER = os.Getenv("ICINGA2_API_USER")
var ICINGA2_API_PASSWORD = os.Getenv("ICINGA2_API_PASSWORD")
var ICINGA2_API_URL = os.Getenv("ICINGA2_API_URL")
var ICINGA2_INSECURE_SKIP_TLS_VERIFY, _ = strconv.ParseBool(os.Getenv("ICINGA2_INSECURE_SKIP_TLS_VERIFY"))

var Icinga2_Server = Server{ICINGA2_API_USER, ICINGA2_API_PASSWORD, ICINGA2_API_URL, ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0, nil}

func TestConnect(t *testing.T) {

	v := os.Getenv("ICINGA2_API_URL")
	if v == "" {
		t.Fatal("ICINGA2_API_URL must be set for acceptance tests")
	}

	v = os.Getenv("ICINGA2_API_USER")
	if v == "" {
		t.Fatal("ICINGA2_API_USER must be set for acceptance tests")
	}

	v = os.Getenv("ICINGA2_API_PASSWORD")
	if v == "" {
		t.Fatal("ICINGA2_API_PASSWORD must be set for acceptance tests")
	}

	var Icinga2_Server = Server{"icinga-test", "icinga", ICINGA2_API_URL, ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0, nil}
	Icinga2_Server.Connect()

	if Icinga2_Server.httpClient == nil {
		t.Errorf("Failed to succesfully connect to Icinga Server")
	}
}

func TestConnectWithBadCredential(t *testing.T) {

	var Icinga2_Server = Server{"unknownUser", "unknownPW", ICINGA2_API_URL, ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0, nil}
	err := Icinga2_Server.Connect()
	if err != nil {
		t.Errorf("Did not fail with bad credentials : %s", err)
	}
}

func TestConnectServerBadURINoVersion(t *testing.T) {

	var Icinga2_Server = Server{ICINGA2_API_USER, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665", ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0, nil}
	result, _ := Icinga2_Server.NewAPIRequest("GET", "/status", nil)

	if result.Code != 404 {
		t.Errorf("Error : Did not get expected 404 error connection to bad URI, with no version.")
	}
}

func TestNewAPIRequest(t *testing.T) {

	result, _ := Icinga2_Server.NewAPIRequest("GET", "/status", nil)

	if result.Code != 200 {
		t.Errorf("%s", result.Status)
	}
}

func TestNewAPIRequestWhileReloading(t *testing.T) {
	mockTransport := httpmock.NewMockTransport()
	mockTransport.RegisterResponder("GET", "https://127.0.0.1:5665/status",
		httpmock.ResponderFromMultipleResponses(
			[]*http.Response{
				httpmock.NewStringResponse(http.StatusServiceUnavailable, `{"status":"Icinga is reloading"}`),
				httpmock.NewStringResponse(http.StatusOK, `{}`),
			},
			t.Log),
	)

	tries := 0
	server := Server{ICINGA2_API_USER, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665", ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", tries, 0, nil}
	server.createHttpClient()
	server.httpClient.Transport = mockTransport

	results, err := server.NewAPIRequest("GET", "/status", nil)

	if err == nil {
		t.Errorf("expected error 'icinga is reloading', got nil")
		return
	}

	if results.Code != http.StatusServiceUnavailable {
		t.Errorf("expected code %d, got %d", http.StatusServiceUnavailable, results.Code)
	}

	if err.Error() != "icinga is reloading" {
		t.Errorf("expected error 'icinga is reloading', got '%v'", err)
	}

	info := mockTransport.GetCallCountInfo()
	calls, ok := info["GET https://127.0.0.1:5665/status"]

	if !ok {
		t.Errorf("cannot find mock stats in %v", info)
		return
	}

	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestNewAPIRequestWhileReloadingWithRetries(t *testing.T) {
	mockTransport := httpmock.NewMockTransport()
	mockTransport.RegisterResponder("GET", "https://127.0.0.1:5665/status",
		httpmock.ResponderFromMultipleResponses(
			[]*http.Response{
				httpmock.NewStringResponse(http.StatusServiceUnavailable, `{"status":"Icinga is reloading"}`),
				httpmock.NewStringResponse(http.StatusOK, `{}`),
			},
			t.Log),
	)

	tries := 2
	server := Server{ICINGA2_API_USER, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665", ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", tries, 0, nil}
	server.createHttpClient()
	server.httpClient.Transport = mockTransport

	results, err := server.NewAPIRequest("GET", "/status", nil)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if results.Code != http.StatusOK {
		t.Errorf("expected code %d, got %d", http.StatusOK, results.Code)
	}

	info := mockTransport.GetCallCountInfo()
	calls, ok := info["GET https://127.0.0.1:5665/status"]

	if !ok {
		t.Errorf("cannot find mock stats in %v", info)
		return
	}

	if calls != tries {
		t.Errorf("expected %d calls, got %d", tries, calls)
	}
}
