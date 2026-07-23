package iapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
)

// ScheduleDowntime Schedule a downtime
// https://icinga.com/docs/icinga-2/latest/doc/12-icinga2-api/#icinga2-api-actions-schedule-downtime
func (server *Server) ScheduleDowntime(ctx context.Context, t string, filter string, author string, comment string, startTime int64, endTime int64, fixed bool, duration int64, allServices bool, triggerName string, childOptions string) ([]string, error) {
	payload := DowntimeScheduleRequest{
		Type:         t,
		Filter:       filter,
		Author:       author,
		Comment:      comment,
		StartTime:    startTime,
		EndTime:      endTime,
		Fixed:        fixed,
		Duration:     duration,
		AllServices:  allServices,
		TriggerName:  triggerName,
		ChildOptions: childOptions,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal downtime payload: %v", err)
	}

	r, err := server.NewAPIRequest(ctx, "POST", "/actions/schedule-downtime", payloadJSON, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to POST on the API: %v", err)
	}

	if r.Code != http.StatusOK {
		return nil, fmt.Errorf("%d, got %d: %v", http.StatusOK, r.Code, r)
	}

	var results []DowntimeScheduleResponse
	if len(r.Results) > 0 {
		if err := json.Unmarshal(r.Results, &results); err != nil {
			return nil, fmt.Errorf("failed to unmarshal the downtime response: %v", err)
		}
	}

	var names []string
	for _, downtime := range results {
		names = append(names, downtime.Name)
	}
	return names, nil
}

// RemoveDowntime Remove a downtime
// https://icinga.com/docs/icinga-2/latest/doc/12-icinga2-api/#remove-downtime
func (server *Server) RemoveDowntime(ctx context.Context, downtime string, author string) error {
	payload := DowntimeRemoveRequest{
		Downtime: downtime,
		Author:   author,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal downtime payload: %v", err)
	}

	response, err := server.NewAPIRequest(ctx, "POST", "/actions/remove-downtime", payloadJSON, nil)
	if err != nil {
		return fmt.Errorf("failed to POST on the API: %v", err)
	}

	if !slices.Contains([]int{http.StatusOK, http.StatusNotFound}, response.Code) {
		return fmt.Errorf("expected code %d or %d, got %d: %v", http.StatusOK, http.StatusNotFound, response.Code, response)
	}

	var results []DowntimeRemoveResponse
	if len(response.Results) > 0 {
		if err := json.Unmarshal(response.Results, &results); err != nil {
			return fmt.Errorf("failed to unmarshal the downtime response: %v", err)
		}
	}

	for _, result := range results {
		if int(result.Code) != http.StatusOK {
			return fmt.Errorf("failed to delete downtime: %s", result.Status)
		}
	}
	return nil
}
