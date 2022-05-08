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
	hostURL   string
	client    *http.Client
	apiKey    string
	debugHTTP bool
}

func NewClient(apiKey string) *Client {
	c := Client{
		client:  &http.Client{Timeout: 10 * time.Second},
		hostURL: HostURL,
	}

	c.apiKey = apiKey

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
