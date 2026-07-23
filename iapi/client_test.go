package iapi

import (
	"context"
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
var ICINGA2_API_CA_CERT_FILE = os.Getenv("ICINGA2_API_CA_CERT_FILE")

var Icinga2_Server *Server

func TestMain(m *testing.M) {
	Icinga2_Server, _ = New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, ICINGA2_API_URL, ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0)
	os.Exit(m.Run())
}

// Helper function to build a test server or skip the test if not configured.
func getTestServer(t *testing.T) *Server {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}
	return Icinga2_Server
}

func TestConnect(t *testing.T) {
	server := getTestServer(t)
	err := server.Connect(context.Background())
	if err != nil {
		t.Fatalf("Failed to connect to Icinga Server: %v", err)
	}

	if server.httpClient == nil {
		t.Errorf("Failed to successfully connect to Icinga Server")
	}
}

func TestConnectCA(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}
	server, err := New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, ICINGA2_API_URL, false, ICINGA2_API_CA_CERT_FILE, 0, 0)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	err = server.Connect(context.Background())
	if err != nil {
		t.Errorf("Failed to successfully connect to Icinga Server: %s", err)
	}

	if server.httpClient == nil {
		t.Errorf("Failed to successfully connect to Icinga Server")
	}
}

func TestConnectWithBadCredential(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}
	server, err := New("unknownUser", "unknownPW", ICINGA2_API_URL, ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	err = server.Connect(context.Background())
	if err != nil {
		t.Errorf("Did not fail with bad credentials : %s", err)
	}
}

func TestConnectServerBadURINoVersion(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}
	server, err := New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665", ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	result, _ := server.NewAPIRequest(context.Background(), "GET", "/status", nil, nil)

	if result.Code != 404 {
		t.Errorf("Error : Did not get expected 404 error connection to bad URI, with no version.")
	}
}

func TestNewAPIRequest(t *testing.T) {
	server := getTestServer(t)
	result, _ := server.NewAPIRequest(context.Background(), "GET", "/status", nil, nil)

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
	server, err := New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665", ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", tries, 0)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	server.httpClient.Transport = mockTransport

	results, err := server.NewAPIRequest(context.Background(), "GET", "/status", nil, nil)

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
	server, err := New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665", ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", tries, 0)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	server.httpClient.Transport = mockTransport

	results, err := server.NewAPIRequest(context.Background(), "GET", "/status", nil, nil)

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
