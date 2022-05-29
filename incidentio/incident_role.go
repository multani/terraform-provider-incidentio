package incidentio

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
	client  *Client
	urlPart string
}

func (c *Client) IncidentRoles() *IncidentRoles {
	return &IncidentRoles{
		client:  c,
		urlPart: "incident_roles",
	}
}

func (i *IncidentRoles) Get(id string) (*IncidentRoleResponse, error) {
	response := &IncidentRoleResponse{}

	if err := i.client.get(i.urlPart, id, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (i *IncidentRoles) Create(role IncidentRole) (*IncidentRoleResponse, error) {
	response := &IncidentRoleResponse{}

	if err := i.client.create(i.urlPart, role, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (i *IncidentRoles) Update(id string, role IncidentRole) (*IncidentRoleResponse, error) {
	response := &IncidentRoleResponse{}

	if err := i.client.update(i.urlPart, id, role, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (i *IncidentRoles) Delete(id string) error {
	return i.client.delete(i.urlPart, id)
}
