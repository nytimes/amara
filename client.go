package amara // import "github.com/nytimes/amara"

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/sethgrid/pester"
)

type Client struct {
	apiKey              string
	team                string
	endpoint            string
	rateLimitProtection RateLimitProtection
	*pester.Client
}

type RateLimitProtection struct {
	enabled          bool
	triggered        bool
	counter          uint
	MinRetryDuration time.Duration
	MaxRetryDuration time.Duration
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
		RateLimitProtection{
			MinRetryDuration: 5 * time.Second,
			MaxRetryDuration: 20 * time.Minute,
		},
		httpClient,
	}
}

func (c *Client) EnableRateLimitProtection() {
	c.rateLimitProtection.enabled = true
	c.rateLimitProtection.triggered = false
	c.rateLimitProtection.counter = 0
}

func (c *Client) DisableRateLimitProtection() {
	c.rateLimitProtection.enabled = false
	c.rateLimitProtection.triggered = false
}

func (c *Client) SetRateLimitProtection(rlp RateLimitProtection) {
	c.rateLimitProtection = rlp
}

func (c *Client) doRequest(method, url string, body io.Reader) ([]byte, error) {
	fmt.Println("Making amara request " + url)
	if c.rateLimitProtection.triggered {
		return nil, fmt.Errorf("Amara API is currently being rate limited. try again later")
	}

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
		if res.StatusCode == http.StatusTooManyRequests && c.rateLimitProtection.enabled {
			if err := c.rateLimitProtection.start(res); err != nil {
				return nil, err
			}
		}

		// TODO: parse the data
		return nil, fmt.Errorf("status %d: %s", res.StatusCode, data)
	}

	c.rateLimitProtection.reset()

	return data, nil
}

func (r *RateLimitProtection) start(res *http.Response) error {
	var wait time.Duration

	retryAfterStr := res.Header.Get("Retry-After")
	if retryAfterStr != "" {
		retryAfter, err := strconv.Atoi(retryAfterStr)
		if err == nil {
			wait = time.Duration(retryAfter) * time.Second
		} else {
			retryAt, err := time.Parse(time.RFC1123, retryAfterStr)

			wait = time.Until(retryAt)
			if err != nil || wait < 0 {
				return fmt.Errorf("Invalid Retry-After header: %s", retryAfterStr)
			}
		}
	} else {
		wait = time.Duration(math.Pow(2, float64(r.counter))) * r.MinRetryDuration
	}

	fmt.Printf("Wait time: %s\n\n", wait)

	if wait < r.MinRetryDuration {
		wait = r.MinRetryDuration
	}

	if wait > r.MaxRetryDuration {
		wait = r.MaxRetryDuration
	}

	r.triggered = true
	r.counter += 1
	time.AfterFunc(wait, func() {
		r.triggered = false
	})

	return nil
}

func (r *RateLimitProtection) reset() {
	r.triggered = false
	r.counter = 0
}
