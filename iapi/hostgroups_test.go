package iapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestHostgroups(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}
	icingaServer, err := New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, ICINGA2_API_URL, ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		t.Run("Hostgroup", func(t *testing.T) {
			name := "docker-servers"
			displayName := "Docker Host Servers"
			_, err := icingaServer.CreateHostgroup(ctx, name, displayName, "")
			if err != nil {
				t.Error(err)
			}
		})
		t.Run("HostgroupWithZone", func(t *testing.T) {
			name := "docker-servers-zone"
			displayName := "Docker Host Servers on zone"
			_, err := icingaServer.CreateHostgroup(ctx, name, displayName, "master")
			if err != nil {
				t.Error(err)
			}
		})

	})

	t.Run("Read", func(t *testing.T) {
		t.Run("ValidHostgroup", func(t *testing.T) {
			name := "linux-servers"
			_, err := icingaServer.GetHostgroup(ctx, name)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("InvalidHostgroup", func(t *testing.T) {
			name := "irix-servers"
			_, err := icingaServer.GetHostgroup(ctx, name)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("Update", func(t *testing.T) {
		hostgroupName := "someHostgroupName"
		firstDisplayName := "some Hostgroup Display Name"
		_, err := icingaServer.CreateHostgroup(ctx, hostgroupName, firstDisplayName, "")
		if err != nil {
			t.Fatalf("failed to create hostgroup: %v", err)
		}
		defer icingaServer.DeleteHostgroup(ctx, hostgroupName)

		secondDisplayName := "other hostgroup display name"
		attrs := HostgroupAttrs{
			DisplayName: secondDisplayName,
		}
		updatedHostgroup, err := icingaServer.UpdateHostgroup(ctx, hostgroupName, attrs)
		if err != nil {
			t.Fatalf("failed to update hostgroup: %v", err)
		}
		if len(updatedHostgroup) == 0 {
			t.Fatalf("expected non-empty updated hostgroup slice")
		}
		if secondDisplayName != updatedHostgroup[0].Attrs.DisplayName {
			t.Errorf("expected display_name fields to be equal, got different")
		}
	})

	t.Run("Exists", func(t *testing.T) {
		t.Run("HostgroupFound", func(t *testing.T) {
			name := "docker-servers"
			exists, err := icingaServer.HostgroupExists(ctx, name)
			if err != nil {
				t.Error(err)
			}

			if !exists {
				t.Error("host group must exist")
			}
		})

		t.Run("HostgroupNotFound", func(t *testing.T) {
			name := "irix-servers"
			exists, err := icingaServer.HostgroupExists(ctx, name)
			if err != nil {
				t.Error(err)
			}
			if exists {
				t.Error("host group must not exist")
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		// Delete Hostgroup created via API. Should succeed
		t.Run("Hostgroup", func(t *testing.T) {
			name := "docker-servers"
			err := icingaServer.DeleteHostgroup(ctx, name)
			if err != nil {
				t.Error(err)
			}
		})
		t.Run("HostgroupWithZone", func(t *testing.T) {
			name := "docker-servers-zone"
			err := icingaServer.DeleteHostgroup(ctx, name)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("HostgroupNonAPI", func(t *testing.T) {
			name := "linux-servers"
			err := icingaServer.DeleteHostgroup(ctx, name)
			if err == nil {
				t.Error("expected error deleting non-api hostgroup, got nil")
			}
		})

		t.Run("HostgroupDNE", func(t *testing.T) {
			name := "docker-servers"
			err := icingaServer.DeleteHostgroup(ctx, name)
			if err == nil || err.Error() != "No objects found." {
				t.Error(err)
			}
		})
	})
}

func newMockedHostgroupServer(t *testing.T, tries int) (*Server, *httpmock.MockTransport) {
	t.Helper()
	mockTransport := httpmock.NewMockTransport()
	server, err := New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, "https://127.0.0.1:5665", ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", tries, 0)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	server.httpClient.Transport = mockTransport
	return server, mockTransport
}

func hostgroupURL(name string) string {
	return fmt.Sprintf("https://127.0.0.1:5665/objects/hostgroups/%s", url.PathEscape(name))
}

func hostgroupResponse(name string) *http.Response {
	return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"results":[{"type":"HostGroup","name":"%s","attrs":{"display_name":"%s"}}]}`, name, name))
}

func hostgroupNotFoundResponse() *http.Response {
	return httpmock.NewStringResponse(http.StatusNotFound, `{"error":404,"status":"No objects found."}`)
}

func assertCallCount(t *testing.T, mockTransport *httpmock.MockTransport, method, endpoint string, expected int) {
	t.Helper()
	key := fmt.Sprintf("%s %s", method, endpoint)
	if calls := mockTransport.GetCallCountInfo()[key]; calls != expected {
		t.Errorf("expected %d calls to '%s', got %d", expected, key, calls)
	}
}

func TestCreateHostgroupWithDeadlineExceeded(t *testing.T) {
	name := "go-icinga2-api-hostgroup-1"
	server, mockTransport := newMockedHostgroupServer(t, 0)
	mockTransport.RegisterResponder(http.MethodPut, hostgroupURL(name), httpmock.NewErrorResponder(context.DeadlineExceeded))
	mockTransport.RegisterResponder(http.MethodGet, hostgroupURL(name), func(*http.Request) (*http.Response, error) {
		return hostgroupNotFoundResponse(), nil
	})

	if _, err := server.CreateHostgroup(context.Background(), name, name, ""); err == nil {
		t.Errorf("expected an error when the host group is never created, got nil")
	}
}

func TestCreateHostgroupWithDeadlineExceededAndRetries(t *testing.T) {
	name := "go-icinga2-api-hostgroup-2"
	server, mockTransport := newMockedHostgroupServer(t, 2)
	mockTransport.RegisterResponder(http.MethodPut, hostgroupURL(name), httpmock.NewErrorResponder(context.DeadlineExceeded))
	mockTransport.RegisterResponder(http.MethodGet, hostgroupURL(name), httpmock.ResponderFromMultipleResponses(
		[]*http.Response{hostgroupNotFoundResponse(), hostgroupResponse(name)},
		t.Log,
	))

	hostgroups, err := server.CreateHostgroup(context.Background(), name, name, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hostgroups) != 1 || hostgroups[0].Name != name {
		t.Errorf("expected host group '%s', got %v", name, hostgroups)
	}

	assertCallCount(t, mockTransport, http.MethodPut, hostgroupURL(name), 1)
	assertCallCount(t, mockTransport, http.MethodGet, hostgroupURL(name), 2)
}

func TestDeleteHostgroupWithDeadlineExceeded(t *testing.T) {
	name := "go-icinga2-api-hostgroup-3"
	server, mockTransport := newMockedHostgroupServer(t, 0)
	mockTransport.RegisterResponder(http.MethodDelete, hostgroupURL(name), httpmock.NewErrorResponder(context.DeadlineExceeded))
	mockTransport.RegisterResponder(http.MethodGet, hostgroupURL(name), func(*http.Request) (*http.Response, error) {
		return hostgroupResponse(name), nil
	})

	if err := server.DeleteHostgroup(context.Background(), name); err == nil {
		t.Errorf("expected an error when the host group is never deleted, got nil")
	}
}

func TestDeleteHostgroupWithDeadlineExceededAndRetries(t *testing.T) {
	name := "go-icinga2-api-hostgroup-4"
	server, mockTransport := newMockedHostgroupServer(t, 2)
	mockTransport.RegisterResponder(http.MethodDelete, hostgroupURL(name), httpmock.NewErrorResponder(context.DeadlineExceeded))
	mockTransport.RegisterResponder(http.MethodGet, hostgroupURL(name), httpmock.ResponderFromMultipleResponses(
		[]*http.Response{hostgroupResponse(name), hostgroupNotFoundResponse()},
		t.Log,
	))

	if err := server.DeleteHostgroup(context.Background(), name); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertCallCount(t, mockTransport, http.MethodDelete, hostgroupURL(name), 1)
	assertCallCount(t, mockTransport, http.MethodGet, hostgroupURL(name), 2)
}
