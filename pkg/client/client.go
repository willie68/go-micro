package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/willie68/go-micro/internal/logging"
	"github.com/willie68/go-micro/internal/serror"
	"github.com/willie68/go-micro/pkg/pmodel"
)

// Client the client to be used for other parts
type Client struct {
	url      string
	tenant   string
	clt      http.Client
	ctx      context.Context
	insecure bool
}

// NewClient creating a new client for the go-micro service template
func NewClient(u, t string) (*Client, error) {
	cl := Client{
		tenant: t,
	}
	err := cl.init(u)
	if err != nil {
		return nil, err
	}
	return &cl, nil
}

func (c *Client) init(u string) error {
	timeout := time.Second * 5
	c.insecure = false
	ul, err := url.Parse(u)
	if err != nil {
		return err
	}
	if ul.Hostname() == "127.0.0.1" {
		c.insecure = true
		timeout = time.Second * 360
	}
	c.url = fmt.Sprintf("%s/api/v1", u)
	c.ctx = context.Background()

	tns := &http.Transport{
		// #nosec G402 -- fine for internal traffic
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.insecure,
		},
	}

	c.clt = http.Client{
		Timeout:   timeout,
		Transport: tns,
	}
	return nil
}

// GetAddresses getting all config names
func (c *Client) GetAddresses() (*[]pmodel.Address, error) {
	res, err := c.Get("addresses")
	if err != nil {
		logging.Root.Error(fmt.Sprintf("get request failed: %v", err))
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		logging.Root.Error(fmt.Sprintf("get bad response: %d", res.StatusCode))
		return nil, ReadErr(res)
	}
	var l []pmodel.Address
	err = ReadJSON(res, &l)
	if err != nil {
		logging.Root.Error(fmt.Sprintf("parsing response failed: %v", err))
		return nil, err
	}
	return &l, nil
}

// GetAddress getting the address of a id
func (c *Client) GetAddress(n string) (*pmodel.Address, error) {
	res, err := c.Get(fmt.Sprintf("addresses/%s", n))
	if err != nil {
		logging.Root.Error(fmt.Sprintf("get request failed: %v", err))
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		logging.Root.Error(fmt.Sprintf("get bad response: %d", res.StatusCode))
		return nil, ReadErr(res)
	}
	var cd pmodel.Address
	err = ReadJSON(res, &cd)
	if err != nil {
		logging.Root.Error(fmt.Sprintf("parsing response failed: %v", err))
		return nil, err
	}
	return &cd, nil
}

// CreateAddress create the address
func (c *Client) CreateAddress(pcd pmodel.Address) (string, error) {
	res, err := c.PostJSON("addresses/", pcd)
	if err != nil {
		logging.Root.Error(fmt.Sprintf("put request failed: %v", err))
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		logging.Root.Error(fmt.Sprintf("put bad response: %d", res.StatusCode))
		return "", ReadErr(res)
	}
	cd := struct {
		ID string `json:"id"`
	}{}
	err = ReadJSON(res, &cd)
	if err != nil {
		logging.Root.Error(fmt.Sprintf("parsing response failed: %v", err))
		return "", err
	}
	return cd.ID, nil
}

// DeleteAddress getting the address
func (c *Client) DeleteAddress(n string) (bool, error) {
	res, err := c.Delete(fmt.Sprintf("addresses/%s", n))
	if err != nil {
		logging.Root.Error(fmt.Sprintf("delete request failed: %v", err))
		return false, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		if res.StatusCode != http.StatusNotFound {
			logging.Root.Error(fmt.Sprintf("delete bad response: %d", res.StatusCode))
			return false, ReadErr(res)
		}
		return false, nil
	}
	return true, nil
}

// Get getting something from the endpoint
func (c *Client) Get(endpoint string) (*http.Response, error) {
	req, err := c.newRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

// Post posting something to the endpoint
func (c *Client) Post(endpoint, contentType string, body io.Reader) (*http.Response, error) {
	req, err := c.newRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.do(req)
}

// PostJSON posting a json string to the endpoint
func (c *Client) PostJSON(endpoint string, body any) (*http.Response, error) {
	byt, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return c.Post(endpoint, "application/json", bytes.NewBuffer(byt))
}

// Delete sending a delete to an endpoint
func (c *Client) Delete(endpoint string) (*http.Response, error) {
	req, err := c.newRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

// do a request with logging
func (c *Client) do(req *http.Request) (*http.Response, error) {
	ul := req.URL.RequestURI()
	res, err := c.clt.Do(req)
	if err != nil {
		logging.Root.Error(fmt.Sprintf("request %s %s error: %v", req.Method, ul, err))
	} else {
		logging.Root.Info(fmt.Sprintf("request %s %s returned %s", req.Method, ul, res.Status))
	}
	return res, err
}

func (c *Client) newRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	ul := fmt.Sprintf("%s/%s", c.url, endpoint)
	logging.Root.Debug(fmt.Sprintf("creating request %s %s", method, ul))
	req, err := http.NewRequestWithContext(c.ctx, method, ul, body)
	if err != nil {
		logging.Root.Error(fmt.Sprintf("cannot create request %s %s", method, ul))
		return nil, err
	}
	req.Header.Set("tenant", c.tenant)
	logging.Root.Debug(fmt.Sprintf("request %s %s", method, ul))
	return req, nil
}

// ReadJSON read the given response as json
func ReadJSON(res *http.Response, dst any) error {
	return json.NewDecoder(res.Body).Decode(&dst)
}

// ReadErr read the given response as an error
func ReadErr(res *http.Response) error {
	var serr serror.Serr
	err := ReadJSON(res, &serr)
	if err != nil {
		byt, _ := io.ReadAll(res.Body)
		return serror.New(res.StatusCode, "bad-response", string(byt))
	}
	return &serr
}
