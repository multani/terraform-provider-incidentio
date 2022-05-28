package incidentio_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/multani/terraform-provider-incidentio/incidentio"
)

func TestCustomFieldsGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/custom_fields/01G44T2BWJY0ZMV945X32RAJ5C", r.URL.String())
		require.Equal(t, "GET", r.Method)
		_, err := w.Write([]byte(`
		{
			"custom_field": {
				"id": "01G44T2BWJY0ZMV945X32RAJ5C",
				"name": "Affected Team",
				"description": "The team which was responsible for resolving this incident.",
				"field_type": "multi_select",
				"required": "always",
				"show_before_creation": false,
				"show_before_closure": true,
				"options": [],
				"created_at": "2022-05-28T07:46:07.385Z",
				"updated_at": "2022-05-28T07:46:07.385Z"
			}
		}
		  `))
		require.NoError(t, err)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := incidentio.NewClient("foobar").WithHostURL(server.URL)

	field, err := client.CustomFields().Get("01G44T2BWJY0ZMV945X32RAJ5C")
	require.NoError(t, err)

	assert.Equal(t, "Affected Team", field.CustomField.Name)
	assert.Equal(t, "The team which was responsible for resolving this incident.", field.CustomField.Description)
	assert.Equal(t, incidentio.FieldRequirement("always"), field.CustomField.Required)
	assert.Equal(t, false, field.CustomField.ShowBeforeCreation)
	assert.Equal(t, true, field.CustomField.ShowBeforeClosure)
	assert.Equal(t, []incidentio.CustomFieldOption{}, field.CustomField.Options)
	assert.Equal(t, incidentio.FieldType("multi_select"), field.CustomField.FieldType)
}

func TestCustomFieldsCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		require.Equal(t, r.URL.String(), "/v1/custom_fields")
		require.Equal(t, r.Method, "POST")

		field := &incidentio.CustomField{}
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &field)
		require.NoError(t, err)

		// Send response
		w.WriteHeader(http.StatusCreated)
		fieldResp := &incidentio.CustomFieldResponse{
			CustomField: incidentio.CustomFieldMetadata{
				Id:          "id123",
				CreatedAt:   "",
				UpdatedAt:   "",
				CustomField: *field,
			},
		}
		body, err = json.Marshal(fieldResp)
		require.NoError(t, err)

		_, err = w.Write(body)
		require.NoError(t, err)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := incidentio.NewClient("foobar").WithHostURL(server.URL)

	request := incidentio.CustomField{
		Name:               "some name",
		Description:        "some description",
		Required:           "always",
		ShowBeforeCreation: false,
		ShowBeforeClosure:  true,
		FieldType:          "number",
	}

	response, err := client.CustomFields().Create(request)
	require.NoError(t, err)

	assert.Equal(t, "id123", response.CustomField.Id)
	assert.Equal(t, request.Name, response.CustomField.Name)
	assert.Equal(t, request.Description, response.CustomField.Description)
}

func TestCustomFieldsUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		require.Equal(t, r.URL.String(), "/v1/custom_fields/id123")
		require.Equal(t, r.Method, "PUT")

		field := &incidentio.CustomField{}
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &field)
		require.NoError(t, err)

		// Send response
		w.WriteHeader(http.StatusOK)
		fieldResp := &incidentio.CustomFieldResponse{
			CustomField: incidentio.CustomFieldMetadata{
				Id:          "id123",
				CreatedAt:   "",
				UpdatedAt:   "",
				CustomField: *field,
			},
		}
		body, err = json.Marshal(fieldResp)
		require.NoError(t, err)

		_, err = w.Write(body)
		require.NoError(t, err)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := incidentio.NewClient("foobar").WithHostURL(server.URL)

	request := incidentio.CustomField{
		Name:               "some name",
		Description:        "some description",
		Required:           "always",
		ShowBeforeCreation: false,
		ShowBeforeClosure:  true,
		FieldType:          "number",
	}

	response, err := client.CustomFields().Update("id123", request)
	require.NoError(t, err)

	assert.Equal(t, "id123", response.CustomField.Id)
	assert.Equal(t, request.Name, response.CustomField.Name)
	assert.Equal(t, request.Description, response.CustomField.Description)
}

func TestCustomFieldsDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		require.Equal(t, r.URL.String(), "/v1/custom_fields/id123")
		require.Equal(t, r.Method, "DELETE")

		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		require.Empty(t, body)

		// Send response
		w.WriteHeader(http.StatusNoContent)
		_, err = w.Write([]byte(""))
		require.NoError(t, err)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := incidentio.NewClient("foobar").WithHostURL(server.URL)

	err := client.CustomFields().Delete("id123")
	require.NoError(t, err)
}
