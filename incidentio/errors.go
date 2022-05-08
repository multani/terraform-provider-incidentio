package incidentio

import (
	"encoding/json"
	"strings"
)

type IncidentIOErrorResponse struct {
	Type      string            `json:"type"`
	Status    int               `json:"status"`
	RequestID string            `json:"request_id"`
	Errors    []IncidentIOError `json:"errors"`
}

type IncidentIOError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Source  SourceError `json:"source"`
}

type SourceError struct {
	Field   string `json:"field"`
	Pointer string `json:"pointer"`
}

func NewErrors(body []byte) error {
	errorResponse := &IncidentIOErrorResponse{}
	err := json.Unmarshal(body, &errorResponse)
	if err != nil {
		return err
	}

	return errorResponse

}

func (e *IncidentIOErrorResponse) Error() string {
	var builder strings.Builder

	builder.WriteString(e.Type + ": ")

	for _, err := range e.Errors {
		builder.WriteString(err.Code + ":" + err.Message)
	}

	return builder.String()
}
