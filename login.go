package mcapi

func (c *Client) Login(userID, password string) error {
	body := map[string]interface{}{
		"user_id":  userID,
		"password": password,
	}

	var result struct {
		Data struct {
			APIKey string `json:"apikey"`
		} `json:"data"`
	}

	if err := c.post(&result, body, "login"); err != nil {
		return err
	}

	c.APIKey = result.Data.APIKey
	return nil
}
