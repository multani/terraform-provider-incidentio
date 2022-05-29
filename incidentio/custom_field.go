package incidentio

import (
	"fmt"
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
	client  *Client
	urlPart string
}

func (c *Client) CustomFields() *CustomFields {
	return &CustomFields{
		client:  c,
		urlPart: "custom_fields",
	}
}

func (i *CustomFields) Get(id string) (*CustomFieldResponse, error) {
	response := &CustomFieldResponse{}

	if err := i.client.get(i.urlPart, id, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (i *CustomFields) Create(field CustomField) (*CustomFieldResponse, error) {
	response := &CustomFieldResponse{}

	if err := i.client.create(i.urlPart, field, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (i *CustomFields) Update(id string, field CustomField) (*CustomFieldResponse, error) {
	response := &CustomFieldResponse{}

	if err := i.client.update(i.urlPart, id, field, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (i *CustomFields) Delete(id string) error {
	return i.client.delete(i.urlPart, id)
}
