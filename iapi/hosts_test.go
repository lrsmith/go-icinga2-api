package iapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGetValidHost(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	hostname := "c1-mysql-1"

	_, err := Icinga2_Server.GetHost(context.Background(), hostname)

	if err != nil {
		t.Error(err)
	}
}

func TestGetInvalidHost(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	hostname := "c2-mysql-1"
	_, err := Icinga2_Server.GetHost(context.Background(), hostname)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateSimpleHost(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	hostname := "go-icinga2-api-1"
	IPAddress := "127.0.0.2"
	CheckCommand := "hostalive"

	_, err := Icinga2_Server.CreateHost(context.Background(), hostname, IPAddress, "", CheckCommand, nil, nil, nil, "")

	if err != nil {
		t.Error(err)
	}
}

func TestCreateSimpleIPv6Host(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	hostname := "go-icinga2-api-3"
	IPAddress := "127.0.0.2"
	IP6Address := "::1"
	CheckCommand := "hostalive"

	_, err := Icinga2_Server.CreateHost(context.Background(), hostname, IPAddress, IP6Address, CheckCommand, nil, nil, nil, "")

	if err != nil {
		t.Error(err)
	}

	// Delete host after creating it.
	deleteErr := Icinga2_Server.DeleteHost(context.Background(), hostname)
	if deleteErr != nil {
		t.Error(deleteErr)
	}
}

func TestCreateHostWithVariables(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	hostname := "go-icinga2-api-2"
	IPAddress := "127.0.0.3"
	CheckCommand := "hostalive"

	variables := make(map[string]interface{})

	variables["os"] = "Linux"
	variables["creator"] = "Terraform"
	variables["urls"] = []string{"test-url1.example.com", "test-url2.example.com"}

	hosts, err := Icinga2_Server.CreateHost(context.Background(), hostname, IPAddress, "", CheckCommand, variables, nil, nil, "")
	if err != nil {
		t.Error(err)
	}

	vars := hostVars(t, hosts, hostname)
	if vars["os"] != "Linux" {
		t.Errorf("expected variable os to be Linux, got %v", vars["os"])
	}
	if vars["creator"] != "Terraform" {
		t.Errorf("expected variable creator to be Terraform, got %v", vars["creator"])
	}

	// Delete host after creating it.
	deleteErr := Icinga2_Server.DeleteHost(context.Background(), hostname)
	if deleteErr != nil {
		t.Error(deleteErr)
	}
}

func TestCreateHostWithTemplates(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	hostname := "go-icinga2-api-2"
	IPAddress := "127.0.0.3"
	CheckCommand := "hostalive"

	templates := []string{"go-icinga2-api-test-linux"}

	hosts, err := Icinga2_Server.CreateHost(context.Background(), hostname, IPAddress, "", CheckCommand, nil, templates, nil, "")
	if err != nil {
		t.Error(err)
	}

	// Icinga resolves the imports of the imported templates too and reports the
	// host itself first
	expected := []string{hostname, "go-icinga2-api-test-linux", "generic-host"}
	for _, host := range hosts {
		if host.Name != hostname {
			continue
		}
		if !slices.Equal(host.Attrs.Templates, expected) {
			t.Errorf("expected templates %v, got %v", expected, host.Attrs.Templates)
		}
	}

	// A variable of the template only reaches the host when the template is
	// imported instead of being stored as an attribute
	if vars := hostVars(t, hosts, hostname); vars["os"] != "Linux" {
		t.Errorf("expected variable os to be inherited from the template, got %v", vars["os"])
	}

	// Delete host after creating it.
	deleteErr := Icinga2_Server.DeleteHost(context.Background(), hostname)
	if deleteErr != nil {
		t.Error(deleteErr)
	}
}

func TestCreateHostWithTemplatesAndVariables(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	hostname := "go-icinga2-api-2"
	IPAddress := "127.0.0.3"
	CheckCommand := "hostalive"

	templates := []string{"go-icinga2-api-test-linux"}
	variables := map[string]interface{}{"role": "database"}

	hosts, err := Icinga2_Server.CreateHost(context.Background(), hostname, IPAddress, "", CheckCommand, variables, templates, nil, "")
	if err != nil {
		t.Error(err)
	}

	vars := hostVars(t, hosts, hostname)
	if vars["role"] != "database" {
		t.Errorf("expected variable role to be database, got %v", vars["role"])
	}
	if vars["os"] != "Linux" {
		t.Errorf("expected variable os to be inherited from the template, got %v", vars["os"])
	}

	// Delete host after creating it.
	deleteErr := Icinga2_Server.DeleteHost(context.Background(), hostname)
	if deleteErr != nil {
		t.Error(deleteErr)
	}
}

func TestCreateHostWithUnknownTemplate(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	hostname := "go-icinga2-api-unknown-template"
	IPAddress := "127.0.0.3"
	CheckCommand := "hostalive"

	templates := []string{"go-icinga2-api-test-does-not-exist"}

	_, err := Icinga2_Server.CreateHost(context.Background(), hostname, IPAddress, "", CheckCommand, nil, templates, nil, "")
	if err == nil {
		t.Fatalf("expected the creation of %s to be rejected", hostname)
	}

	if !strings.Contains(err.Error(), "Import references unknown template") {
		t.Errorf("expected the reason of the rejection, got %s", err)
	}
}

func TestCreateHostWithGroup(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	hostname := "go-icinga2-api-2"
	IPAddress := "127.0.0.3"
	CheckCommand := "hostalive"
	Group := []string{"linux-servers"}

	_, err := Icinga2_Server.CreateHost(context.Background(), hostname, IPAddress, "", CheckCommand, nil, nil, Group, "")
	if err != nil {
		t.Error(err)
	}

	// Delete host after creating it.
	deleteErr := Icinga2_Server.DeleteHost(context.Background(), hostname)
	if deleteErr != nil {
		t.Error(deleteErr)
	}
}

func TestCreateHostWithZone(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	hostname := "go-icinga2-api-2"
	IPAddress := "127.0.0.3"
	CheckCommand := "hostalive"
	Group := []string{"linux-servers"}

	_, err := Icinga2_Server.CreateHost(context.Background(), hostname, IPAddress, "", CheckCommand, nil, nil, Group, "master")
	if err != nil {
		t.Error(err)
	}

	// Delete host after creating it.
	deleteErr := Icinga2_Server.DeleteHost(context.Background(), hostname)
	if deleteErr != nil {
		t.Error(deleteErr)
	}
}

func TestCreateHostWithDeadlineExceeded(t *testing.T) {
	hostname := "go-icinga2-api-5"
	IPAddress := "127.0.0.2"
	CheckCommand := "hostalive"

	// Mock for CreateHost
	// Always return context deadline exceeded
	mockTransport := httpmock.NewMockTransport()
	mockTransport.RegisterResponder(
		http.MethodPut,
		fmt.Sprintf("https://127.0.0.1:5665/objects/hosts/%s", url.PathEscape(hostname)),
		httpmock.NewErrorResponder(context.DeadlineExceeded),
	)

	server, err := New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665", ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	server.httpClient.Transport = mockTransport

	_, err = server.CreateHost(context.Background(), hostname, IPAddress, "", CheckCommand, nil, nil, nil, "")

	if err == nil {
		t.Errorf("expected context deadline exceeded error, got nil")
	}
}

func TestCreateHostWithDeadlineExceededAndRetries(t *testing.T) {
	hostname := "go-icinga2-api-5"
	IPAddress := "127.0.0.2"
	CheckCommand := "hostalive"

	mockTransport := httpmock.NewMockTransport()

	// Mock for CreateHost
	// Return context deadline exceeded
	mockTransport.RegisterResponder(http.MethodPut,
		fmt.Sprintf("https://127.0.0.1:5665/objects/hosts/%s", url.PathEscape(hostname)),
		httpmock.NewErrorResponder(context.DeadlineExceeded),
	)

	// Mock for GetHost
	// Return 404 to simulate eventual consistency
	// Then 200 to ensure the host has been created successfully
	mockTransport.RegisterResponder(http.MethodGet,
		fmt.Sprintf("https://127.0.0.1:5665/objects/hosts/%s", url.PathEscape(hostname)),
		httpmock.ResponderFromMultipleResponses(
			[]*http.Response{
				httpmock.NewStringResponse(http.StatusNotFound, `{"error":404,"status":"No objects found."}`),
				httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"results":[{"type":"Host","name":"%s","address":"%s","check_command":"%s"}]}`, hostname, IPAddress, CheckCommand)),
			},
			t.Log),
	)

	tries := 2
	server, err := New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665", ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", tries, 0)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	server.httpClient.Transport = mockTransport

	_, err = server.CreateHost(context.Background(), hostname, IPAddress, "", CheckCommand, nil, nil, nil, "")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	info := mockTransport.GetCallCountInfo()

	createHostEndpoint := fmt.Sprintf("%s https://127.0.0.1:5665/objects/hosts/%s", http.MethodPut, url.PathEscape(hostname))
	createCalls, ok := info[createHostEndpoint]

	if !ok {
		t.Errorf("cannot find mock stats for '%s' endpoint in %v", createHostEndpoint, info)
		return
	}

	if createCalls != 1 {
		t.Errorf("expected 1 call to CreateHost, got %d", createCalls)
	}

	getHostEndpoint := fmt.Sprintf("%s https://127.0.0.1:5665/objects/hosts/%s", http.MethodGet, url.PathEscape(hostname))
	getCalls, ok := info[getHostEndpoint]
	if !ok {
		t.Errorf("cannot find mock stats for '%s' endpoint in %v", getHostEndpoint, info)
		return
	}

	if getCalls != 2 {
		t.Errorf("expected 2 calls to GetHost, got %d", getCalls)
	}
}

func TestDeleteHost(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	hostname := "go-icinga2-api-1"

	err := Icinga2_Server.DeleteHost(context.Background(), hostname)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteHostDNE(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	hostname := "go-icinga2-api-1"
	err := Icinga2_Server.DeleteHost(context.Background(), hostname)
	if err == nil || err.Error() != "No objects found." {
		t.Error(err)
	}
}

func TestDeleteHostWithDeadlineExceeded(t *testing.T) {
	hostname := "go-icinga2-api-6"

	// Mock for DeleteHost
	// Always return context deadline exceeded
	mockTransport := httpmock.NewMockTransport()
	mockTransport.RegisterResponder(
		http.MethodDelete,
		fmt.Sprintf("https://127.0.0.1:5665/objects/hosts/%s", url.PathEscape(hostname)),
		httpmock.NewErrorResponder(context.DeadlineExceeded),
	)

	server, err := New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665", ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	server.httpClient.Transport = mockTransport

	err = server.DeleteHost(context.Background(), hostname)
	if err == nil {
		t.Errorf("expected context deadline exceeded error, got nil")
	}
}

func TestDeleteHostWithDeadlineExceededAndRetries(t *testing.T) {
	hostname := "go-icinga2-api-6"
	IPAddress := "127.0.0.2"
	CheckCommand := "hostalive"

	mockTransport := httpmock.NewMockTransport()

	// Mock for DeleteHost
	// Return context deadline exceeded
	mockTransport.RegisterResponder(http.MethodDelete,
		fmt.Sprintf("https://127.0.0.1:5665/objects/hosts/%s", url.PathEscape(hostname)),
		httpmock.NewErrorResponder(context.DeadlineExceeded),
	)

	// Mock for GetHost (via HostExists)
	// Return 200 to simulate eventual consistency
	// Then 404 to ensure the host has been deleted successfully
	mockTransport.RegisterResponder(http.MethodGet,
		fmt.Sprintf("https://127.0.0.1:5665/objects/hosts/%s", url.PathEscape(hostname)),
		httpmock.ResponderFromMultipleResponses(
			[]*http.Response{
				httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"results":[{"type":"Host","name":"%s","address":"%s","check_command":"%s"}]}`, hostname, IPAddress, CheckCommand)),
				httpmock.NewStringResponse(http.StatusNotFound, `{"error":404,"status":"No objects found."}`),
			},
			t.Log),
	)

	tries := 2
	server, err := New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665", ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", tries, 0)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	server.httpClient.Transport = mockTransport

	err = server.DeleteHost(context.Background(), hostname)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	info := mockTransport.GetCallCountInfo()

	deleteHostEndpoint := fmt.Sprintf("%s https://127.0.0.1:5665/objects/hosts/%s", http.MethodDelete, url.PathEscape(hostname))
	deleteCalls, ok := info[deleteHostEndpoint]

	if !ok {
		t.Errorf("cannot find mock stats for '%s' endpoint in %v", deleteHostEndpoint, info)
		return
	}

	if deleteCalls != 1 {
		t.Errorf("expected 1 call to DeleteHost, got %d", deleteCalls)
	}

	getHostEndpoint := fmt.Sprintf("%s https://127.0.0.1:5665/objects/hosts/%s", http.MethodGet, url.PathEscape(hostname))
	getCalls, ok := info[getHostEndpoint]
	if !ok {
		t.Errorf("cannot find mock stats for '%s' endpoint in %v", getHostEndpoint, info)
		return
	}

	if getCalls != 2 {
		t.Errorf("expected 2 calls to GetHost, got %d", getCalls)
	}
}

func TestHostExistsFound(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	hostname := "go-icinga2-api-4"
	IPAddress := "127.0.0.4"
	CheckCommand := "hostalive"

	_, err := Icinga2_Server.CreateHost(context.Background(), hostname, IPAddress, "", CheckCommand, nil, nil, nil, "")

	if err != nil {
		t.Error(err)
	}

	exists, err := Icinga2_Server.HostExists(context.Background(), hostname)
	if err != nil {
		t.Error(err)
	}

	if !exists {
		t.Error("host must exist")
	}

	err = Icinga2_Server.DeleteHost(context.Background(), hostname)
	if err != nil {
		t.Error(err)
	}
}

func TestHostExistsNotFound(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}

	hostname := "go-icinga2-api-4-not-found"

	exists, err := Icinga2_Server.HostExists(context.Background(), hostname)
	if err != nil {
		t.Error(err)
	}

	if exists {
		t.Error("host must not exist")
	}
}

// hostVars returns the resolved variables of a host.
func hostVars(t *testing.T, hosts []HostStruct, hostname string) map[string]interface{} {
	t.Helper()

	for _, host := range hosts {
		if host.Name != hostname {
			continue
		}
		vars, ok := host.Attrs.Vars.(map[string]interface{})
		if !ok {
			t.Fatalf("expected variables on %s, got %v", hostname, host.Attrs.Vars)
		}
		return vars
	}

	t.Fatalf("host %s not found", hostname)
	return nil
}
