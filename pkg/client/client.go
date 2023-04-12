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

// GetConfigs getting all config names
func (c *Client) GetConfigs() (*[]string, error) {
	res, err := c.Get("config")
	if err != nil {
		logging.Logger.Errorf("get request failed: %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		logging.Logger.Errorf("get bad response: %d", res.StatusCode)
		return nil, ReadErr(res)
	}
	var l []string
	err = ReadJSON(res, &l)
	if err != nil {
		logging.Logger.Errorf("parsing response failed: %v", err)
		return nil, err
	}
	return &l, nil
}

// GetConfig getting the config description of a name
func (c *Client) GetConfig(n string) (*pmodel.ConfigDescription, error) {
	res, err := c.Get(fmt.Sprintf("config/%s", n))
	if err != nil {
		logging.Logger.Errorf("get request failed: %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		logging.Logger.Errorf("get bad response: %d", res.StatusCode)
		return nil, ReadErr(res)
	}
	var cd pmodel.ConfigDescription
	err = ReadJSON(res, &cd)
	if err != nil {
		logging.Logger.Errorf("parsing response failed: %v", err)
		return nil, err
	}
	return &cd, nil
}

// GetMyConfig getting the config size of me
func (c *Client) GetMyConfig() (*pmodel.ConfigDescription, error) {
	res, err := c.Get("config/_own")
	if err != nil {
		logging.Logger.Errorf("get request failed: %v", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		logging.Logger.Errorf("get bad response: %d", res.StatusCode)
		return nil, ReadErr(res)
	}
	var cd pmodel.ConfigDescription
	err = ReadJSON(res, &cd)
	if err != nil {
		logging.Logger.Errorf("parsing response failed: %v", err)
		return nil, err
	}
	return &cd, nil
}

// PutConfig getting the config description of a name
func (c *Client) PutConfig(pcd pmodel.ConfigDescription) (string, error) {
	res, err := c.PostJSON("config/", pcd)
	if err != nil {
		logging.Logger.Errorf("put request failed: %v", err)
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		logging.Logger.Errorf("put bad response: %d", res.StatusCode)
		return "", ReadErr(res)
	}
	cd := struct {
		ID string `json:"id"`
	}{}
	err = ReadJSON(res, &cd)
	if err != nil {
		logging.Logger.Errorf("parsing response failed: %v", err)
		return "", err
	}
	return cd.ID, nil
}

// DeleteConfig getting the config description of a name
func (c *Client) DeleteConfig(n string) (bool, error) {
	res, err := c.Delete(fmt.Sprintf("config/%s", n))
	if err != nil {
		logging.Logger.Errorf("delete request failed: %v", err)
		return false, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		if res.StatusCode != http.StatusNotFound {

			logging.Logger.Errorf("delete bad response: %d", res.StatusCode)
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
		logging.Logger.Errorf("request %s %s error: %v", req.Method, ul, err)
	} else {
		logging.Logger.Infof("request %s %s returned %s", req.Method, ul, res.Status)
	}
	return res, err
}

func (c *Client) newRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	ul := fmt.Sprintf("%s/%s", c.url, endpoint)
	logging.Logger.Debugf("creating request %s %s", method, ul)
	req, err := http.NewRequestWithContext(c.ctx, method, ul, body)
	if err != nil {
		logging.Logger.Errorf("cannot create request %s %s", method, ul)
		return nil, err
	}
	req.Header.Set("tenant", c.tenant)
	logging.Logger.Debugf("request %s %s", method, ul)
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
