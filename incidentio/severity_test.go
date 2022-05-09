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

func TestSeveritiesGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/severities/01FCNDV6P870EA6S7TK1DSYDG0", r.URL.String())
		require.Equal(t, "GET", r.Method)
		_, err := w.Write([]byte(`
		{
			"severity": 
			{
				"created_at": "2021-08-17T13:28:57.801578Z",
				"description": "It's not really that bad, everyone chill",
				"id": "01FCNDV6P870EA6S7TK1DSYDG0",
				"name": "Minor",
				"rank": 1,
				"updated_at": "2021-08-17T13:28:57.801578Z"
			}
		}
		  `))
		require.NoError(t, err)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := incidentio.NewClient("foobar").WithHostURL(server.URL)

	response, err := client.Severities().Get("01FCNDV6P870EA6S7TK1DSYDG0")
	require.NoError(t, err)

	assert.Equal(t, "Minor", response.Severity.Name)
	assert.Equal(t, "It's not really that bad, everyone chill", response.Severity.Description)
	assert.Equal(t, int64(1), response.Severity.Rank)
}

func TestSeveritiesCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		require.Equal(t, r.URL.String(), "/v1/severities")
		require.Equal(t, r.Method, "POST")

		severity := &incidentio.Severity{}
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &severity)
		require.NoError(t, err)

		// Send response
		w.WriteHeader(http.StatusCreated)
		response := &incidentio.SeverityResponse{
			Severity: incidentio.SeverityMetadata{
				Id:        "id123",
				CreatedAt: "",
				UpdatedAt: "",
				Severity:  *severity,
			},
		}
		body, err = json.Marshal(response)
		require.NoError(t, err)

		_, err = w.Write(body)
		require.NoError(t, err)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := incidentio.NewClient("foobar").WithHostURL(server.URL)

	request := incidentio.Severity{
		Name:        "some name",
		Description: "some description",
		Rank:        42,
	}

	response, err := client.Severities().Create(request)
	require.NoError(t, err)

	assert.Equal(t, "id123", response.Severity.Id)
	assert.Equal(t, request.Name, response.Severity.Name)
	assert.Equal(t, request.Description, response.Severity.Description)
}

func TestSeveritiesUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		require.Equal(t, r.URL.String(), "/v1/severities/id123")
		require.Equal(t, r.Method, "PUT")

		severity := &incidentio.Severity{}
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &severity)
		require.NoError(t, err)

		// Send response
		w.WriteHeader(http.StatusOK)
		response := &incidentio.SeverityResponse{
			Severity: incidentio.SeverityMetadata{
				Id:        "id123",
				CreatedAt: "",
				UpdatedAt: "",
				Severity:  *severity,
			},
		}
		body, err = json.Marshal(response)
		require.NoError(t, err)

		_, err = w.Write(body)
		require.NoError(t, err)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := incidentio.NewClient("foobar").WithHostURL(server.URL)

	request := incidentio.Severity{
		Name:        "some name",
		Description: "some description",
		Rank:        64,
	}

	response, err := client.Severities().Update("id123", request)
	require.NoError(t, err)

	assert.Equal(t, "id123", response.Severity.Id)
	assert.Equal(t, request.Name, response.Severity.Name)
	assert.Equal(t, request.Description, response.Severity.Description)
	assert.Equal(t, int64(64), response.Severity.Rank)
}

func TestSeveritiesDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		require.Equal(t, r.URL.String(), "/v1/severities/id123")
		require.Equal(t, r.Method, "DELETE")

		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		require.Empty(t, body)

		// Send response
		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write([]byte(""))
		require.NoError(t, err)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := incidentio.NewClient("foobar").WithHostURL(server.URL)

	err := client.Severities().Delete("id123")
	require.NoError(t, err)
}
