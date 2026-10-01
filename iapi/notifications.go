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

// notificationCreateRequest is the payload of a notification creation.
// Attributes are addressed by path, which is what makes "vars.os" merge with the
// variables inherited from the templates where a whole "vars" dictionary would
// replace them.
type notificationCreateRequest struct {
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Attrs     map[string]interface{} `json:"attrs"`
	Templates []string               `json:"templates,omitempty"`
}

// CreateNotification creates a notification.
// The keys of vars are variable names, without the "vars." prefix.
func (server *Server) CreateNotification(ctx context.Context, name, hostname, command, servicename string, interval int, users []string, vars map[string]string, templates []string) ([]NotificationStruct, error) {
	// Only send the attributes the caller provided, an empty one would override
	// the value inherited from the templates. Icinga rejects an empty command.
	attrs := make(map[string]interface{})
	if command != "" {
		attrs["command"] = command
	}
	if servicename != "" {
		attrs["service_name"] = servicename
	}
	if interval != 0 {
		attrs["interval"] = interval
	}
	if users != nil {
		attrs["users"] = users
	}

	// Addressing each variable merges them with the ones inherited from the
	// templates, assigning the whole "vars" dictionary would replace them.
	for variable, value := range vars {
		attrs["vars."+variable] = value
	}

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(notificationCreateRequest{
		Name:  name,
		Type:  "Notification",
		Attrs: attrs,
		// Templates are imports rather than an attribute. Icinga stores an
		// "attrs.templates" array on the object verbatim and imports nothing.
		Templates: templates,
	})
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

	return nil, objectCreateError(results)
}

// UpdateNotification updates a Notification with its attrs in-place
func (server *Server) UpdateNotification(ctx context.Context, name string, attrs NotificationAttrs) ([]NotificationStruct, error) {
	// Icinga only accepts "service_name" and "templates" when the notification
	// is created and rejects the whole update when either is sent along.
	payload := make(map[string]interface{})
	if attrs.Command != "" {
		payload["command"] = attrs.Command
	}
	if attrs.Users != nil {
		payload["users"] = attrs.Users
	}
	if attrs.Interval != 0 {
		payload["interval"] = attrs.Interval
	}
	if attrs.Vars != nil {
		payload["vars"] = attrs.Vars
	}

	body, err := json.Marshal(map[string]interface{}{"attrs": payload})
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

	if updateErr := objectUpdateError(r); updateErr != nil {
		return nil, updateErr
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
