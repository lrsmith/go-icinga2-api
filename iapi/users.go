package iapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetUser ...
func (server *Server) GetUser(ctx context.Context, name string) ([]UserStruct, error) {
	var users []UserStruct
	_, err := server.NewAPIRequest(ctx, "GET", "/objects/users/"+name, nil, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// CreateUser ...
func (server *Server) CreateUser(ctx context.Context, name, email string, variables map[string]string) ([]UserStruct, error) {
	var newAttrs UserAttrs
	newAttrs.Email = email
	newAttrs.Vars = variables

	var newUser UserStruct
	newUser.Name = name
	newUser.Type = "User"
	newUser.Attrs = newAttrs

	// Create JSON from completed struct
	payloadJSON, marshalErr := json.Marshal(newUser)
	if marshalErr != nil {
		return nil, marshalErr
	}

	// Make the API request to create the users.
	results, err := server.NewAPIRequest(ctx, "PUT", "/objects/users/"+name, payloadJSON, nil)
	if err != nil {
		return nil, err
	}

	if results.Code == 200 {
		return server.GetUser(ctx, name)
	}

	return nil, fmt.Errorf("%s", results.ErrorString)
}

// UpdateUser updates a User with its attrs in-place
func (server *Server) UpdateUser(ctx context.Context, name string, attrs UserAttrs) ([]UserStruct, error) {
	user := UserStruct{
		Attrs: attrs,
	}

	body, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	r, err := server.NewAPIRequest(ctx, "POST", "/objects/users/"+name, body, nil)
	if err != nil {
		return nil, err
	}

	if r.Code != http.StatusOK {
		return nil, fmt.Errorf("expected %d, got %d", http.StatusOK, r.Code)
	}

	return server.GetUser(ctx, name)
}

// DeleteUser ...
func (server *Server) DeleteUser(ctx context.Context, name string) error {
	results, err := server.NewAPIRequest(ctx, "DELETE", "/objects/users/"+name+"?cascade=1", nil, nil)
	if err != nil {
		return err
	}

	if results.Code == 200 {
		return nil
	}

	return fmt.Errorf("%s", results.ErrorString)
}
