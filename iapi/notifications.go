package iapi

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetNotification ...
func (server *Server) GetNotification(ctx context.Context, name string) ([]NotificationStruct, error) {
	var notifications []NotificationStruct
	_, err := server.NewAPIRequest(ctx, "GET", "/objects/notifications/"+name, nil, &notifications)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

// CreateNotification ...
func (server *Server) CreateNotification(ctx context.Context, name, hostname, command, servicename string, interval int, users []string, vars map[string]string, templates []string) ([]NotificationStruct, error) {
	var newAttrs NotificationAttrs
	newAttrs.Command = command
	newAttrs.Users = users
	newAttrs.Servicename = servicename
	newAttrs.Interval = interval
	newAttrs.Vars = vars
	newAttrs.Templates = templates

	var newNotification NotificationStruct
	newNotification.Name = name
	newNotification.Type = "Notification"
	newNotification.Attrs = newAttrs

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(newNotification)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// Make the API request to create the notification.
	results, err := server.NewAPIRequest(ctx, "PUT", "/objects/notifications/"+name, payloadJSON, nil)
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		return server.GetNotification(ctx, name)
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
}

// UpdateNotification updates a Notification with its attrs in-place
func (server *Server) UpdateNotification(ctx context.Context, name string, attrs NotificationAttrs) ([]NotificationStruct, error) {
	notification := NotificationStruct{
		Attrs: attrs,
	}

	body, err := json.Marshal(notification)
	if err != nil {
		return nil, err
	}

	r, err := server.NewAPIRequest(ctx, "POST", "/objects/notifications/"+name, body, nil)
	if err != nil {
		return nil, err
	}

	// Accept 200 OK
	if r.Code != 200 {
		return nil, fmt.Errorf("expected 200, got %d: %s", r.Code, r.ErrorString)
	}

	return server.GetNotification(ctx, name)
}

// DeleteNotification ...
func (server *Server) DeleteNotification(ctx context.Context, name string) error {
	results, err := server.NewAPIRequest(ctx, "DELETE", "/objects/notifications/"+name+"?cascade=1", nil, nil)
	if err != nil {
		return err
	}

	if results.Code == 200 {
		return nil
	}

	return fmt.Errorf("%s", results.ErrorString)
}
