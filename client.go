package mcapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/materials-commons/gomcapi/pkg/urlpath"
	"gopkg.in/resty.v1"
)

type Client struct {
	APIKey  string
	BaseURL string
}

var ErrAuth = errors.New("authentication")
var ErrMCAPI = errors.New("mcapi")

var tlsConfig = tls.Config{InsecureSkipVerify: true}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: urlpath.Join(baseURL, "v3"),
	}
}

func (c *Client) r() *resty.Request {
	return resty.SetTLSClientConfig(&tlsConfig).R()
}

func (c *Client) Join(paths ...string) string {
	return urlpath.Join(c.BaseURL, paths...)
}

func (c *Client) post(result, body interface{}, paths ...string) error {
	resp, err := r().SetResult(&result).SetBody(body).Post(c.Join(paths...))
	return c.getAPIError(resp, err)
}

func (c *Client) getAPIError(resp *resty.Response, err error) error {
	switch {
	case err != nil:
		return err
	case resp.RawResponse.StatusCode == 401:
		return ErrAuth
	case resp.RawResponse.StatusCode > 299:
		return c.toErrorFromResponse(resp)
	default:
		return nil
	}
}

func (c *Client) toErrorFromResponse(resp *resty.Response) error {
	var er struct {
		Data struct {
			Error string `json:"error"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body(), &er); err != nil {
		return errors.WithMessage(ErrMCAPI, fmt.Sprintf("(HTTP Status: %d)- unable to parse json error response: %s", resp.RawResponse.StatusCode, err))
	}

	return errors.WithMessage(ErrMCAPI, fmt.Sprintf("(HTTP Status: %d)- %s", resp.RawResponse.StatusCode, er.Data.Error))
}
