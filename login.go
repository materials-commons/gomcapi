package mcapi

import "github.com/materials-commons/gomcapi/pkg/urlpath"

func Login(userID, password, url string) (*Client, error) {
	c := NewConnection(url, "")
	body := map[string]interface{}{
		"user_id":  userID,
		"password": password,
	}

	var res struct {
		Data struct {
			APIKey string `json:"apikey"`
		} `json:"data"`
	}

	resp, err := r().SetResult(&res).SetBody(body).Post(urlpath.Join(url, "v3", "login"))

	if err := c.getAPIError(resp, err); err != nil {
		return nil, err
	}

	c.APIKey = res.Data.APIKey
	return c, nil
}
