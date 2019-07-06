package webpowerswitch

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	user = "admin"
	pass = "incubus1"
)

// Outlet is info about the current state of an outlet on the switch
type Outlet struct {
	Name           string `json:"name"`
	Critical       bool   `json:"critical"`
	TransientState bool   `json:"transient_state"`
	CycleDelay     string `json:"cycle_delay"`
	PhysicalState  bool   `json:"physical_state"`
	Locked         bool   `json:"locked"`
	State          bool   `json:"state"`
}

// Client is a http client for interating with a web power switch
type Client struct {
	BaseURL            *url.URL
	httpClient         *http.Client
	username, password string
}

// NewClient creates a new client
func NewClient(baseURL, username, password string) (*Client, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	cli := Client{
		BaseURL:  base,
		username: username,
		password: password,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	return &cli, nil
}

// GetOutlets returns all outlets on he switch
func (c *Client) GetOutlets() ([]*Outlet, error) {
	rel := &url.URL{Path: "restapi/relay/outlets/"}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.username, c.password)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Request failed, got status " + strconv.Itoa(resp.StatusCode) + " from web power switch")
	}

	var outlets []*Outlet
	err = json.NewDecoder(resp.Body).Decode(&outlets)
	return outlets, err
}

// TurnOutletOn returns all outlets on he switch
func (c *Client) TurnOutletOn(outletID string) error {
	rel := &url.URL{Path: "restapi/relay/outlets/" + outletID + "/state/"}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("PUT", u.String(), bytes.NewBuffer([]byte("value=true")))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-CSRF", "x")
	req.SetBasicAuth(c.username, c.password)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errors.New("Request failed, got status " + strconv.Itoa(resp.StatusCode) + " from web power switch")
	}

	return err
}

// TurnOutletOff turns off an outlet given an id
func (c *Client) TurnOutletOff(outletID string) error {
	rel := &url.URL{Path: "restapi/relay/outlets/" + outletID + "/state/"}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("PUT", u.String(), bytes.NewBuffer([]byte("value=false")))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-CSRF", "x")
	req.SetBasicAuth(c.username, c.password)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errors.New("Request failed, got status " + strconv.Itoa(resp.StatusCode) + " from web power switch")
	}

	return err
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
