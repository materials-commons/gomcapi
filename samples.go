package mcapi

func (c *Client) CreateSample() (*Sample, error) {
	var result struct {
		Data Sample `json:"data"`
	}

	body := struct{}{}

	if err := c.post(&result, body, "createSample"); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (c *Client) AddSampleToProcess() (*Sample, error) {
	var result struct {
		Data Sample `json:"data"`
	}

	body := struct{}{}

	if err := c.post(&result, body, "addSampleToProcess"); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
