package incidentio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type FieldType string

const (
	SingleSelect FieldType = "single_select"
	MultiSelect  FieldType = "multi_select"
	Text         FieldType = "text"
	Link         FieldType = "link"
	Numeric      FieldType = "numeric"
)

func ParseFieldType(s string) (*FieldType, error) {
	v := FieldType(s)

	switch v {
	case SingleSelect, MultiSelect, Text, Link, Numeric:
		return &v, nil
	}

	return nil, fmt.Errorf("%v is not a valid field type", s)
}

type FieldRequirement string

const (
	Never         FieldRequirement = "never"
	BeforeClosure FieldRequirement = "before_closure"
	Always        FieldRequirement = "always"
)

func ParseFieldRequirement(s string) (*FieldRequirement, error) {
	v := FieldRequirement(s)

	switch v {
	case Never, BeforeClosure, Always:
		return &v, nil
	}

	return nil, fmt.Errorf("%v is not a valid field requirement", s)
}

type CustomField struct {
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	Required           FieldRequirement    `json:"required"`
	ShowBeforeClosure  bool                `json:"show_before_closure"`
	ShowBeforeCreation bool                `json:"show_before_creation"`
	FieldType          FieldType           `json:"field_type"`
	Options            []CustomFieldOption `json:"options"`
}

type CustomFieldMetadata struct {
	CustomField

	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CustomFieldResponse struct {
	CustomField CustomFieldMetadata `json:"custom_field"`
}

// CustomFields is used to query custom field options
type CustomFields struct {
	client *Client
}

func (c *Client) CustomFields() *CustomFields {
	return &CustomFields{client: c}
}

func (i *CustomFields) Get(id string) (*CustomFieldResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("You must specify the ID of the custom field option to get")
	}

	request, err := i.client.newRequest("GET", fmt.Sprintf("/v1/custom_fields/%s", id), nil)
	if err != nil {
		return nil, err
	}

	res, body, err := i.client.doRequest(request)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, NewErrors(body)
	}

	response := &CustomFieldResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *CustomFields) Create(role CustomField) (*CustomFieldResponse, error) {
	rb, err := json.Marshal(role)
	if err != nil {
		return nil, err
	}

	request, err := i.client.newRequest("POST", "/v1/custom_fields", strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	res, body, err := i.client.doRequest(request)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		return nil, NewErrors(body)
	}

	response := &CustomFieldResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *CustomFields) Update(id string, role CustomField) (*CustomFieldResponse, error) {
	rb, err := json.Marshal(role)
	if err != nil {
		return nil, err
	}

	request, err := i.client.newRequest("PUT", fmt.Sprintf("/v1/custom_fields/%s", id), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	res, body, err := i.client.doRequest(request)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, NewErrors(body)
	}

	response := &CustomFieldResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *CustomFields) Delete(id string) error {
	request, err := i.client.newRequest("DELETE", fmt.Sprintf("/v1/custom_fields/%s", id), nil)
	if err != nil {
		return err
	}

	res, body, err := i.client.doRequest(request)

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return err
}
