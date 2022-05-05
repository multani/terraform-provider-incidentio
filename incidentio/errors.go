package incidentio

import "encoding/json"

type IncidentIOError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	// {"source":{"field":"","pointer":"shortform"}}
}

type IncidentIOErrorResponse struct {
	Type      string            `json:"type"`
	Status    int               `json:"status"`
	RequestID string            `json:"request_id"`
	Errors    []IncidentIOError `json:"errors"`
}

func NewErrors(body []byte) (*IncidentIOErrorResponse, error) {
	response := &IncidentIOErrorResponse{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil

}

func (e *IncidentIOErrorResponse) Error() string {
	return e.Type
}
