package mcapi

import (
	"crypto/tls"
	"errors"

	"github.com/materials-commons/gomcapi/pkg/urlpath"
	"gopkg.in/resty.v1"
)

type Client struct {
	APIKey  string
	BaseURL string
}

var ErrBadConfig = errors.New("bad configuration")

var tlsConfig = tls.Config{InsecureSkipVerify: true}

func NewConnection(BaseURL, APIKey string) *Client {
	return &Client{
		APIKey:  APIKey,
		BaseURL: urlpath.Join(BaseURL, "v3"),
	}
}

func (c *Client) r() *resty.Request {
	return resty.SetTLSClientConfig(&tlsConfig).R()
}

func Login(userID, password, url string) (*Client, error) {
	c := NewConnection(url, "")
	body := map[string]interface{}{
		"user_id":  userID,
		"password": password,
	}
	var resp struct {
		APIKey string `json:"apikey"`
	}
	_, err := r().SetResult(&resp).SetBody(body).Post(urlpath.Join(url, "v3", "login"))
	if err != nil {
		return nil, err
	}
	c.APIKey = resp.APIKey
	return c, nil
}
