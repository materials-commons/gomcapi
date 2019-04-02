package mcapi

func (c *Client) CreateSample(projectID, experimentID, name string, attributes []Property) (*Sample, error) {
	var result struct {
		Data Sample `json:"data"`
	}

	if attributes == nil {
		attributes = make([]Property, 0)
	}

	body := struct {
		ProjectID    string     `json:"project_id"`
		ExperimentID string     `json:"experiment_id"`
		Name         string     `json:"name"`
		Attributes   []Property `json:"attributes"`
	}{
		ProjectID:    projectID,
		ExperimentID: experimentID,
		Name:         name,
		Attributes:   attributes,
	}

	if err := c.post(&result, body, "createSample"); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

type ConnectSampleToProcess struct {
	ProcessID     string
	SampleID      string
	PropertySetID string
	Transform     bool
}

func (c *Client) AddSampleToProcess(projectID, experimentID string, connect ConnectSampleToProcess) (*Sample, error) {
	var result struct {
		Data Sample `json:"data"`
	}

	body := struct {
		ProjectID     string `json:"project_id"`
		ExperimentID  string `json:"experiment_id"`
		ProcessID     string `json:"process_id"`
		SampleID      string `json:"sample_id"`
		PropertySetID string `json:"property_set_id"`
		Transform     bool   `json:"transform"`
	}{
		ProjectID:     projectID,
		ExperimentID:  experimentID,
		ProcessID:     connect.ProcessID,
		SampleID:      connect.SampleID,
		PropertySetID: connect.PropertySetID,
		Transform:     connect.Transform,
	}

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
