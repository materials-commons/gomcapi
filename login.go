package mcapi

import "github.com/materials-commons/gomcapi/pkg/urlpath"

func (c *Client) Login(userID, password string) error {
	body := map[string]interface{}{
		"user_id":  userID,
		"password": password,
	}

	var res struct {
		Data struct {
			APIKey string `json:"apikey"`
		} `json:"data"`
	}

	resp, err := r().SetResult(&res).SetBody(body).Post(urlpath.Join(c.BaseURL, "login"))

	if err := c.getAPIError(resp, err); err != nil {
		return err
	}

	c.APIKey = res.Data.APIKey
	return nil
}
