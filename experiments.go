package mcapi

func (c *Client) CreateExperiment(projectID, name, description string) (*Experiment, error) {
	var result struct {
		Data Experiment `json:"data"`
	}

	body := map[string]interface{}{
		"project_id":  projectID,
		"name":        name,
		"description": description,
	}

	if err := c.post(&result, body, "createExperimentInProject"); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
