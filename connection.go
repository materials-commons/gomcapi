package mcapi

import (
	"crypto/tls"
	"errors"

	"github.com/materials-commons/config"
	"github.com/materials-commons/gomcapi/pkg/urlpath"
	"gopkg.in/resty.v1"
)

type Conn struct {
	APIKey  string
	BaseURL string
}

var ErrBadConfig = errors.New("bad configuration")

var tlsConfig = tls.Config{InsecureSkipVerify: true}

func NewConnection(BaseURL, APIKey string) *Conn {
	return &Conn{
		APIKey:  APIKey,
		BaseURL: urlpath.Join(BaseURL, "v3"),
	}
}

func ConnectionFromDefaultConfig() (*Conn, error) {
	apikey := config.GetString("apikey")
	baseURL := config.GetString("mcurl")
	if apikey == "" || baseURL == "" {
		return nil, ErrBadConfig
	}
	return &Conn{APIKey: apikey, BaseURL: urlpath.Join(baseURL, "v3")}, nil
}

func (c *Conn) r() *resty.Request {
	return resty.SetTLSClientConfig(&tlsConfig).R()
}

func Login(userID, password, url string) (*Conn, error) {
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

func LoginUsingDefaultConfig(userID, password string) (*Conn, error) {
	baseURL := config.GetString("mcurl")
	if baseURL == "" {
		return nil, ErrBadConfig
	}

	return Login(userID, password, baseURL)
}
