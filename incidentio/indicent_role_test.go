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

func TestIncidentRolesGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/incident_roles/01FCNDV6P870EA6S7TK1DSYDG0", r.URL.String())
		require.Equal(t, "GET", r.Method)
		_, err := w.Write([]byte(`
		{
			"incident_role": {
			  "created_at": "2021-08-17T13:28:57.801578Z",
			  "description": "The person currently coordinating the incident",
			  "id": "01FCNDV6P870EA6S7TK1DSYDG0",
			  "instructions": "Take point on the incident; Make sure people are clear on responsibilities",
			  "name": "Incident Lead",
			  "required": true,
			  "role_type": "lead",
			  "shortform": "lead",
			  "updated_at": "2021-08-17T13:28:57.801578Z"
			}
		}
		  `))
		require.NoError(t, err)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := incidentio.NewClient("foobar").WithHostURL(server.URL)

	role, err := client.IncidentRoles().Get("01FCNDV6P870EA6S7TK1DSYDG0")
	require.NoError(t, err)

	assert.Equal(t, "Incident Lead", role.IncidentRole.Name)
	assert.Equal(t, "The person currently coordinating the incident", role.IncidentRole.Description)
	assert.Equal(t, true, role.IncidentRole.Required)
	assert.Equal(t, "Take point on the incident; Make sure people are clear on responsibilities", role.IncidentRole.Instructions)
	assert.Equal(t, "lead", role.IncidentRole.ShortForm)
}

func TestIncidentRolesCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		require.Equal(t, r.URL.String(), "/v1/incident_roles")
		require.Equal(t, r.Method, "POST")

		role := &incidentio.IncidentRole{}
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &role)
		require.NoError(t, err)

		// Send response
		w.WriteHeader(http.StatusCreated)
		roleResp := &incidentio.IncidentRoleResponse{
			IncidentRole: incidentio.IncidentRoleMetadata{
				Id:           "id123",
				CreatedAt:    "",
				UpdatedAt:    "",
				RoleType:     "foo",
				IncidentRole: *role,
			},
		}
		body, err = json.Marshal(roleResp)
		require.NoError(t, err)

		_, err = w.Write(body)
		require.NoError(t, err)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := incidentio.NewClient("foobar").WithHostURL(server.URL)

	request := incidentio.IncidentRole{
		Name:         "some name",
		Description:  "some description",
		Required:     true,
		Instructions: "some instructions",
		ShortForm:    "some short form",
	}

	response, err := client.IncidentRoles().Create(request)
	require.NoError(t, err)

	assert.Equal(t, "id123", response.IncidentRole.Id)
	assert.Equal(t, request.Name, response.IncidentRole.Name)
	assert.Equal(t, request.Description, response.IncidentRole.Description)
}

func TestIncidentRolesUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		require.Equal(t, r.URL.String(), "/v1/incident_roles/id123")
		require.Equal(t, r.Method, "PUT")

		role := &incidentio.IncidentRole{}
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &role)
		require.NoError(t, err)

		// Send response
		w.WriteHeader(http.StatusOK)
		roleResp := &incidentio.IncidentRoleResponse{
			IncidentRole: incidentio.IncidentRoleMetadata{
				Id:           "id123",
				CreatedAt:    "",
				UpdatedAt:    "",
				RoleType:     "foo",
				IncidentRole: *role,
			},
		}
		body, err = json.Marshal(roleResp)
		require.NoError(t, err)

		_, err = w.Write(body)
		require.NoError(t, err)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := incidentio.NewClient("foobar").WithHostURL(server.URL)

	request := incidentio.IncidentRole{
		Name:         "some name",
		Description:  "some description",
		Required:     true,
		Instructions: "some instructions",
		ShortForm:    "some short form",
	}

	response, err := client.IncidentRoles().Update("id123", request)
	require.NoError(t, err)

	assert.Equal(t, "id123", response.IncidentRole.Id)
	assert.Equal(t, request.Name, response.IncidentRole.Name)
	assert.Equal(t, request.Description, response.IncidentRole.Description)
}

func TestIncidentRolesDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		require.Equal(t, r.URL.String(), "/v1/incident_roles/id123")
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

	err := client.IncidentRoles().Delete("id123")
	require.NoError(t, err)
}
