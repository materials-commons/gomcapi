package mcapi

func (c *Client) CreateSample(projectID, experimentID, name string, attributes []Property) (*Sample, error) {
	var result struct {
		Data Sample `json:"data"`
	}

	body := struct {
		ProjectID    string     `json:"project_id"`
		ExperimentID string     `json:"experiment_id"`
		Name         string     `json:"name"`
		Attributes   []Property `json:"attributes"`
	}{}

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

func (c *Client) AddMeasurementsToSampleInProcess() (*Sample, error) {
	var result struct {
		Data Sample `json:"data"`
	}

	body := struct{}{}

	if err := c.post(&result, body, "addMeasurementsToSampleInProcess"); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
