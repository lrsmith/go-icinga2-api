package iapi

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestDowntimes(t *testing.T) {
	if ICINGA2_API_URL == "" {
		t.Skip("ICINGA2_API_URL must be set for integration tests")
	}
	icingaServer, err := New(ICINGA2_API_USER, ICINGA2_API_PASSWORD, ICINGA2_API_URL, ICINGA2_INSECURE_SKIP_TLS_VERIFY, "", 0, 0)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		t.Run("Downtime", func(t *testing.T) {
			hostname := "go-icinga2-api-dt-create"
			IPAddress := "127.0.0.2"
			CheckCommand := "hostalive"
			_, err := icingaServer.CreateHost(ctx, hostname, IPAddress, "", CheckCommand, nil, nil, nil, "")
			if err != nil {
				t.Error(err)
			}

			dt := "Host"
			filter := fmt.Sprintf("host.name==\"%s\"", hostname)
			author := "author"
			comment := "comment"
			startTime := time.Now().Unix()
			endTime := time.Now().Unix() + 3600
			fixed := true
			var duration int64
			var allServices bool
			var triggerName, childOptions string

			names, err := icingaServer.ScheduleDowntime(ctx, dt, filter, author, comment, startTime, endTime, fixed, duration, allServices, triggerName, childOptions)
			if err != nil {
				t.Error(err)
			}

			if len(names) == 0 {
				t.Error("No downtime found")
			}

			err = icingaServer.DeleteHost(ctx, hostname)
			if err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("Downtime", func(t *testing.T) {
			hostname := "go-icinga2-api-dt-create"
			IPAddress := "127.0.0.2"
			CheckCommand := "hostalive"
			_, err := icingaServer.CreateHost(ctx, hostname, IPAddress, "", CheckCommand, nil, nil, nil, "")
			if err != nil {
				t.Error(err)
			}

			dt := "Host"
			filter := fmt.Sprintf("host.name==\"%s\"", hostname)
			author := "author"
			comment := "comment"
			startTime := time.Now().Unix()
			endTime := time.Now().Unix() + 3600
			fixed := true
			var duration int64
			var allServices bool
			var triggerName, childOptions string

			names, err := icingaServer.ScheduleDowntime(ctx, dt, filter, author, comment, startTime, endTime, fixed, duration, allServices, triggerName, childOptions)
			if err != nil {
				t.Error(err)
			}

			if len(names) == 0 {
				t.Error("No downtime found")
			}

			err = icingaServer.RemoveDowntime(ctx, hostname, author)
			if err != nil {
				t.Error(err)
			}

			err = icingaServer.DeleteHost(ctx, hostname)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("DowntimeAlreadyRemoved", func(t *testing.T) {
			hostname := "go-icinga2-api-dt-create"
			IPAddress := "127.0.0.2"
			CheckCommand := "hostalive"
			_, err := icingaServer.CreateHost(ctx, hostname, IPAddress, "", CheckCommand, nil, nil, nil, "")
			if err != nil {
				t.Error(err)
			}

			err = icingaServer.RemoveDowntime(ctx, hostname, "author")
			if err != nil {
				t.Error(err)
			}

			err = icingaServer.DeleteHost(ctx, hostname)
			if err != nil {
				t.Error(err)
			}
		})
	})
}
