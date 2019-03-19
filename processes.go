package mcapi

func (c *Client) CreateProcess(projectID, experimentID, name string, setups []Setup) (*Process, error) {
	var result struct {
		Data Process `json:"data"`
	}

	body := struct {
		ProjectID    string  `json:"project_id"`
		ExperimentID string  `json:"experiment_id"`
		Name         string  `json:"name"`
		Attributes   []Setup `json:"attributes"`
	}{
		ProjectID:    projectID,
		ExperimentID: experimentID,
		Name:         name,
		Attributes:   setups,
	}

	if err := c.post(&result, body, "createProcess"); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
