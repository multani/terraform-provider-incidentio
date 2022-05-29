package incidentio

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

const HostURL string = "https://api.incident.io"

type Client struct {
	hostURL   string
	client    *http.Client
	apiKey    string
	debugHTTP bool
}

func NewClient(apiKey string) *Client {
	c := Client{
		client:  &http.Client{Timeout: 10 * time.Second},
		hostURL: HostURL,
		apiKey:  apiKey,
	}

	return &c
}

func (c *Client) WithHostURL(url string) *Client {
	c.hostURL = url
	return c
}

func (c *Client) WithDebug(debug bool) *Client {
	c.debugHTTP = debug
	return c
}

func (c *Client) newRequest(method string, path string, body io.Reader) (*http.Request, error) {

	sep := "/"

	if strings.HasSuffix(c.hostURL, "/") || strings.HasPrefix(path, "/") {
		sep = ""
	}

	url := fmt.Sprintf("%s%s%s", c.hostURL, sep, path)
	return http.NewRequest(method, url, body)
}

func (c *Client) doRequest(req *http.Request) (*http.Response, []byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	if c.debugHTTP {
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			panic(err)
		}
		fmt.Printf("### REQUEST:\n%s\n### /REQUEST\n", string(reqDump))
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	if c.debugHTTP {
		respDump, err := httputil.DumpResponse(res, true)
		if err != nil {
			panic(err)
		}
		fmt.Printf("### RESPONSE:\n%s\n### /RESPONSE\n", string(respDump))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, nil, err
	}

	return res, body, err
}

func (c *Client) get(urlPart string, id string, target any) error {
	if id == "" {
		return fmt.Errorf("you must specify an ID to get")
	}

	url := fmt.Sprintf("/v1/%s/%s", urlPart, id)

	request, err := c.newRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, body, err := c.doRequest(request)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return NewErrors(body)
	}

	if err = json.Unmarshal(body, &target); err != nil {
		return err
	}

	return nil
}

func (c *Client) create(urlPart string, input any, target any) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/v1/%s", urlPart)
	reader := strings.NewReader(string(data))

	request, err := c.newRequest("POST", url, reader)
	if err != nil {
		return err
	}

	res, body, err := c.doRequest(request)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return NewErrors(body)
	}

	if err = json.Unmarshal(body, &target); err != nil {
		return err
	}

	return nil
}

func (c *Client) update(urlPart string, id string, input any, target any) error {
	if id == "" {
		return fmt.Errorf("you must specify an ID to update")
	}

	data, err := json.Marshal(input)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/v1/%s/%s", urlPart, id)
	reader := strings.NewReader(string(data))

	request, err := c.newRequest("PUT", url, reader)
	if err != nil {
		return err
	}

	res, body, err := c.doRequest(request)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return NewErrors(body)
	}

	if err = json.Unmarshal(body, &target); err != nil {
		return err
	}

	return nil
}

func (c *Client) delete(urlPart string, id string) error {
	if id == "" {
		return fmt.Errorf("you must specify an ID to delete")
	}

	url := fmt.Sprintf("/v1/%s/%s", urlPart, id)

	request, err := c.newRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	res, body, err := c.doRequest(request)
	if err != nil {
		return err
	}

	// most of the APIs return "204 No Content" upon successful deletion
	// TODO: for some reasons, severities return 204 on successful deletion
	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusAccepted {
		return NewErrors(body)
	}

	return nil
}
