package amara // import "github.com/nytimes/amara"

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/sethgrid/pester"
)

type Client struct {
	apiKey   string
	team     string
	endpoint string
	*pester.Client
}

func NewClient(apiKey, team string) *Client {
	httpClient := pester.New()
	httpClient.Concurrency = 1
	httpClient.MaxRetries = 5
	httpClient.Backoff = pester.ExponentialBackoff
	httpClient.KeepLog = true
	return &Client{
		apiKey,
		team,
		"https://amara.org/api",
		httpClient,
	}
}

func (c *Client) doRequest(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-api-key", c.apiKey)
	req.Header.Set("X-API-FUTURE", "20190619")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		// TODO: parse the data
		return nil, fmt.Errorf("status %d: %s", res.StatusCode, data)
	}

	return data, nil
}
