package incidentio

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
	client  *Client
	urlPart string
}

func (c *Client) CustomFieldOptions() *CustomFieldOptions {
	return &CustomFieldOptions{
		client:  c,
		urlPart: "custom_field_options",
	}
}

func (i *CustomFieldOptions) Get(id string) (*CustomFieldOptionResponse, error) {
	response := &CustomFieldOptionResponse{}

	if err := i.client.get(i.urlPart, id, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (i *CustomFieldOptions) Create(customFieldOption CustomFieldOption) (*CustomFieldOptionResponse, error) {
	response := &CustomFieldOptionResponse{}

	if err := i.client.create(i.urlPart, customFieldOption, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (i *CustomFieldOptions) Update(id string, customFieldOption CustomFieldOption) (*CustomFieldOptionResponse, error) {
	response := &CustomFieldOptionResponse{}

	if err := i.client.update(i.urlPart, id, customFieldOption, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (i *CustomFieldOptions) Delete(id string) error {
	return i.client.delete(i.urlPart, id)
}
