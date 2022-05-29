package incidentio

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
	client  *Client
	urlPart string
}

func (c *Client) Severities() *Severities {
	return &Severities{
		client:  c,
		urlPart: "severities",
	}
}

func (i *Severities) Get(id string) (*SeverityResponse, error) {
	response := SeverityResponse{}

	if err := i.client.get(i.urlPart, id, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (i *Severities) Create(severity Severity) (*SeverityResponse, error) {
	response := &SeverityResponse{}

	if err := i.client.create(i.urlPart, severity, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (i *Severities) Update(id string, severity Severity) (*SeverityResponse, error) {
	response := &SeverityResponse{}

	if err := i.client.update(i.urlPart, id, severity, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (i *Severities) Delete(id string) error {
	return i.client.delete(i.urlPart, id)
}
