package incidentio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type IncidentRole struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Required     bool   `json:"required"`
	Instructions string `json:"instructions"`
	ShortForm    string `json:"shortform"`
}

type IncidentRoleMetadata struct {
	IncidentRole

	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	RoleType  string `json:"role_type"`
}

type IncidentRoleResponse struct {
	IncidentRole IncidentRoleMetadata `json:"incident_role"`
}

// IncidentRoles is used to query incident roles
type IncidentRoles struct {
	client *Client
}

func (c *Client) IncidentRoles() *IncidentRoles {
	return &IncidentRoles{client: c}
}

func (i *IncidentRoles) Get(id string) (*IncidentRoleResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("You must specify the ID of the incident role to search")
	}

	request, err := i.client.newRequest("GET", fmt.Sprintf("/v1/incident_roles/%s", id), nil)
	if err != nil {
		return nil, err
	}

	res, body, err := i.client.doRequest(request)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	response := &IncidentRoleResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *IncidentRoles) Create(role IncidentRole) (*IncidentRoleResponse, error) {
	rb, err := json.Marshal(role)
	if err != nil {
		return nil, err
	}

	request, err := i.client.newRequest("POST", "/v1/incident_roles", strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	res, body, err := i.client.doRequest(request)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		errRep, err := NewErrors(body)
		if err != nil {
			return nil, err
		}
		return nil, errRep
	}

	response := &IncidentRoleResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *IncidentRoles) Update(id string, role IncidentRole) (*IncidentRoleResponse, error) {
	rb, err := json.Marshal(role)
	if err != nil {
		return nil, err
	}

	request, err := i.client.newRequest("PUT", fmt.Sprintf("/v1/incident_roles/%s", id), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	_, body, err := i.client.doRequest(request)
	if err != nil {
		return nil, err
	}

	response := &IncidentRoleResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *IncidentRoles) Delete(id string) error {
	request, err := i.client.newRequest("DELETE", fmt.Sprintf("/v1/incident_roles/%s", id), nil)
	if err != nil {
		return err
	}

	res, body, err := i.client.doRequest(request)

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return err
}
