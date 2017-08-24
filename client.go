package amara

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	username string
	apiKey   string
	Team     string
	endpoint string
	*http.Client
}

func NewClient(username, apiKey, team string) *Client {
	return &Client{
		username,
		apiKey,
		team,
		"https://amara.org/api",
		&http.Client{
			Timeout: time.Second * 15,
		},
	}
}

func (c *Client) doRequest(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-api-username", c.username)
	req.Header.Set("X-api-key", c.apiKey)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	//TODO: check for res.Status to handle http errors, e.g. sending the same video twice will
	// result in a 400
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
