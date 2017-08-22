package amara

import (
	"fmt"
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

type ReqParams struct {
	Method string
	URL    string
	Body   io.Reader
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

func (c *Client) doRequest(params ReqParams) ([]byte, error) {
	fmt.Println(params.URL)
	req, err := http.NewRequest(params.Method, params.URL, params.Body)
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
