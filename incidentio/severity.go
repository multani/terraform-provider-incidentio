package incidentio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Severity struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Rank        int64  `json:"rank"`
}

type SeverityMetadata struct {
	Severity

	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type SeverityResponse struct {
	Severity SeverityMetadata `json:"severity"`
}

// Severities is used to query severities
type Severities struct {
	client *Client
}

func (c *Client) Severities() *Severities {
	return &Severities{client: c}
}

func (i *Severities) Get(id string) (*SeverityResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("You must specify the ID of the severity to search")
	}

	request, err := i.client.newRequest("GET", fmt.Sprintf("/v1/severities/%s", id), nil)
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

	response := &SeverityResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *Severities) Create(severity Severity) (*SeverityResponse, error) {
	rb, err := json.Marshal(severity)
	if err != nil {
		return nil, err
	}

	request, err := i.client.newRequest("POST", "/v1/severities", strings.NewReader(string(rb)))
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

	response := &SeverityResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *Severities) Update(id string, severity Severity) (*SeverityResponse, error) {
	rb, err := json.Marshal(severity)
	if err != nil {
		return nil, err
	}

	request, err := i.client.newRequest("PUT", fmt.Sprintf("/v1/severities/%s", id), strings.NewReader(string(rb)))
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

	response := &SeverityResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (i *Severities) Delete(id string) error {
	request, err := i.client.newRequest("DELETE", fmt.Sprintf("/v1/severities/%s", id), nil)
	if err != nil {
		return err
	}

	res, body, err := i.client.doRequest(request)

	if res.StatusCode != http.StatusAccepted && res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return err
}
