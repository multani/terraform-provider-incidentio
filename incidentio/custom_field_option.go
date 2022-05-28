package incidentio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CustomFieldOption struct {
	CustomFieldId string `json:"custom_field_id"`
	SortKey       int64  `json:"sort_key"`
	Value         string `json:"value"`
}

type CustomFieldOptionMetadata struct {
	CustomFieldOption

	Id string `json:"id"`
}

type CustomFieldOptionResponse struct {
	CustomFieldOption CustomFieldOptionMetadata `json:"custom_field_option"`
}

// CustomFieldOptions is used to query custom field options
type CustomFieldOptions struct {
	client *Client
}

func (c *Client) CustomFieldOptions() *CustomFieldOptions {
	return &CustomFieldOptions{client: c}
}

func (i *CustomFieldOptions) Get(id string) (*CustomFieldOptionResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("You must specify the ID of the custom field option to get")
	}

	request, err := i.client.newRequest("GET", fmt.Sprintf("/v1/custom_field_options/%s", id), nil)
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

	response := &CustomFieldOptionResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *CustomFieldOptions) Create(role CustomFieldOption) (*CustomFieldOptionResponse, error) {
	rb, err := json.Marshal(role)
	if err != nil {
		return nil, err
	}

	request, err := i.client.newRequest("POST", "/v1/custom_field_options", strings.NewReader(string(rb)))
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

	response := &CustomFieldOptionResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *CustomFieldOptions) Update(id string, role CustomFieldOption) (*CustomFieldOptionResponse, error) {
	rb, err := json.Marshal(role)
	if err != nil {
		return nil, err
	}

	request, err := i.client.newRequest("PUT", fmt.Sprintf("/v1/custom_field_options/%s", id), strings.NewReader(string(rb)))
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

	response := &CustomFieldOptionResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *CustomFieldOptions) Delete(id string) error {
	request, err := i.client.newRequest("DELETE", fmt.Sprintf("/v1/custom_field_options/%s", id), nil)
	if err != nil {
		return err
	}

	res, body, err := i.client.doRequest(request)

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return err
}
