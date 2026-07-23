package iapi

import "encoding/json"

/*
Currently to get something working and that can be refactored there is a lot of duplicate and overlapping decleration. In
part this is because when a variable is defined it is set to a default value. This has been problematic with having an attrs
struct that has all the variables. That struct then cannot be used to create the JSON for the create, without modification,
because it would try and set values that are not configurable via the API. i.e. for hosts "LastCheck" So to keep things moving
duplicate or near duplicate defintions of structs are being defined but can be revisted and refactored later and test will
be in place to ensure everything still works.
*/

// ServiceStruct stores service results
type ServiceStruct struct {
	Attrs ServiceAttrs `json:"attrs"`
	Joins struct{}     `json:"joins"`
	//	Meta  struct{}     `json:"meta"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ServiceAttrs struct {
	CheckCommand string      `json:"check_command"`
	Templates    []string    `json:"templates"`
	Vars         interface{} `json:"vars"`
	//	CheckInterval float64       `json:"check_interval"`
	//	DisplayName   string        `json:"display_name"`
	//	Groups        []interface{} `json:"groups"`
	//Name string `json:"name"`
	//	Type string `json:"type"`
}

// CheckcommandStruct is a struct used to store results from an Icinga2 Checkcommand API call.
type CheckcommandStruct struct {
	Name  string            `json:"name"`
	Type  string            `json:"type"`
	Attrs CheckcommandAttrs `json:"attrs"`
	Joins struct{}          `json:"joins"`
	Meta  struct{}          `json:"meta"`
}

type CheckcommandAttrs struct {
	Arguments interface{} `json:"arguments"`
	Command   []string    `json:"command"`
	Templates []string    `json:"templates"`
	//	Env       interface{} `json:"env"`   				// Available to be set but not supported yet
	//	Package   string      `json:"package"`   		// Available to be set but not supported yet
	//	Timeout   float64     `json:"timeout"`   		// Available to be set but not supported yet
	//	Vars      interface{} `json:"vars"`   			// Available to be set but not supported yet
	//	Zone      string      `json:"zone"`   			// Available to be set but not supported yet
}

// HostgroupStruct is a struct used to store results from an Icinga2 HostGroup API Call. The content are also used to generate the JSON for the CreateHostgroup call
type HostgroupStruct struct {
	Name  string         `json:"name"`
	Type  string         `json:"type"`
	Attrs HostgroupAttrs `json:"attrs"`
}

// HostgroupAttrs ...
type HostgroupAttrs struct {
	DisplayName string `json:"display_name,omitempty"`
	Zone        string `json:"zone,omitempty"`
}

// HostStruct is a struct used to store results from an Icinga2 Host API Call. The content are also used to generate the JSON for the CreateHost call
type HostStruct struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Attrs     HostAttrs `json:"attrs"`
	Meta      struct{}  `json:"meta"`
	Joins     struct{}  `json:"stuct"`
	Templates []string  `json:"templates"`
}

// HostAttrs This is struct lists the attributes that can be set during a CreateHost call. The contents of the struct is converted into JSON
type HostAttrs struct {
	ActionURL    string      `json:"action_url"`
	Address      string      `json:"address"`
	Address6     string      `json:"address6"`
	CheckCommand string      `json:"check_command"`
	DisplayName  string      `json:"display_name"`
	Groups       []string    `json:"groups,omitempty"`
	Notes        string      `json:"notes"`
	NotesURL     string      `json:"notes_url"`
	Templates    []string    `json:"templates,omitempty"`
	Vars         interface{} `json:"vars,omitempty"`
	Zone         string      `json:"zone,omitempty"`
}

// APIResult Stores the results from NewApiRequest
type APIResult struct {
	Error       float64 `json:"error"`
	ErrorString string
	Status      string          `json:"status"`
	Code        int             `json:"code"`
	Results     json.RawMessage `json:"results"`
}

// HostgroupUpdateResult stores the API response after updating a Hostgroup
type HostgroupUpdateResult struct {
	Code   float64 `json:"code"`
	Name   string  `json:"name"`
	Status string  `json:"status"`
}

// HostUpdateResult stores the API response after updating a Host
type HostUpdateResult struct {
	Code   float64 `json:"code"`
	Name   string  `json:"name"`
	Status string  `json:"status"`
}

// APIStatus stores the results of an Icinga2 API Status Call
type APIStatus struct {
	Results []struct {
		Name     string   `json:"name"`
		Perfdata []string `json:"perfdata"`
		Status   struct {
			API struct {
				ConnEndpoints       []interface{} `json:"conn_endpoints"`
				Identity            string        `json:"identity"`
				NotConnEndpoints    []interface{} `json:"not_conn_endpoints"`
				NumConnEndpoints    int           `json:"num_conn_endpoints"`
				NumEndpoints        int           `json:"num_endpoints"`
				NumNotConnEndpoints int           `json:"num_not_conn_endpoints"`
				Zones               struct {
					Master struct {
						ClientLogLag int      `json:"client_log_lag"`
						Connected    bool     `json:"connected"`
						Endpoints    []string `json:"endpoints"`
						ParentZone   string   `json:"parent_zone"`
					} `json:"master"`
				} `json:"zones"`
			} `json:"api"`
		} `json:"status"`
	} `json:"results"`
}

// UserStruct is a struct used to store results from an Icinga2 User API Call. The content are also used to generate the JSON for the CreateUser call
type UserStruct struct {
	Name  string    `json:"name"`
	Type  string    `json:"type"`
	Attrs UserAttrs `json:"attrs"`
	Meta  struct{}  `json:"meta"`
	Joins struct{}  `json:"stuct"`
}

// UserAttrs This is struct lists the attributes that can be set during a CreateUser call. The contents of the struct is converted into JSON
type UserAttrs struct {
	Email string      `json:"email"`
	Vars  interface{} `json:"vars"`
}

// NotificationStruct stores notification results
type NotificationStruct struct {
	Attrs NotificationAttrs `json:"attrs"`
	Joins struct{}          `json:"joins"`
	Name  string            `json:"name"`
	Type  string            `json:"type"`
}

type NotificationAttrs struct {
	Command     string      `json:"command"`
	Users       []string    `json:"users"`
	Servicename string      `json:"service_name"`
	Interval    int         `json:"interval"`
	Vars        interface{} `json:"vars"`
	Templates   []string    `json:"templates"`
}

// DowntimeScheduleRequest Create the API request to schedule a downtime
type DowntimeScheduleRequest struct {
	Type         string `json:"type"`
	Filter       string `json:"filter"`
	Author       string `json:"author"`
	Comment      string `json:"comment"`
	StartTime    int64  `json:"start_time"`
	EndTime      int64  `json:"end_time"`
	Fixed        bool   `json:"fixed"`
	Duration     int64  `json:"duration,omitempty"`
	AllServices  bool   `json:"all_services"`
	TriggerName  string `json:"trigger_name,omitempty"`
	ChildOptions string `json:"child_options,omitempty"`
}

// DowntimeScheduleResponse Store response of the API to schedule a downtime
type DowntimeScheduleResponse struct {
	Code     float64 `json:"code"`
	LegacyID float64 `json:"legacy_id"`
	Name     string  `json:"name"`
	Status   string  `json:"status"`
}

// DowntimeRemoveRequest Create the API request to remove a downtime
type DowntimeRemoveRequest struct {
	Downtime string `json:"downtime"`
	Author   string `json:"author"`
}

// DowntimeRemoveResponse Store response of the API to remove a downtime
type DowntimeRemoveResponse struct {
	Code   float64 `json:"code"`
	Status string  `json:"status"`
}

// ApiUserStruct is a struct used to store results from an Icinga2 ApiUser API Call. The content are also used to generate the JSON for the CreateApiUser call
type ApiUserStruct struct {
	Name  string       `json:"name"`
	Type  string       `json:"type"`
	Attrs ApiUserAttrs `json:"attrs"`
}

// ApiUserAttrs ...
type ApiUserAttrs struct {
	Password    string   `json:"password,omitempty"`
	ClientCN    string   `json:"client_cn,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// HostgroupParams defines all available options related to updating a HostGroup.
type HostgroupParams struct {
	DisplayName string
}
