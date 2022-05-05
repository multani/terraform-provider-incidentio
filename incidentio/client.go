package incidentio

import (
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
	HostURL    string
	HTTPClient *http.Client
	ApiKey     string
}

func NewClient(apiKey string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURL,
	}

	c.ApiKey = apiKey

	return &c, nil
}

func (c *Client) newRequest(method string, path string, body io.Reader) (*http.Request, error) {

	sep := "/"

	if strings.HasSuffix(c.HostURL, "/") || strings.HasPrefix(path, "/") {
		sep = ""
	}

	url := fmt.Sprintf("%s%s%s", c.HostURL, sep, path)
	return http.NewRequest(method, url, body)
}

func (c *Client) doRequest(req *http.Request) (*http.Response, []byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))

	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	respDump, err := httputil.DumpResponse(res, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("RESPONSE:\n%s", string(respDump))

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, nil, err
	}

	return res, body, err
}
