package mcapi

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/materials-commons/hydra/pkg/mcdb/mcmodel"
)

// DataWrapper wraps json responses that have a data key before getting to the data, eg
// {"data": {"id": 1,... }}. It is used to get at the underlying data object.
type DataWrapper struct {
	Data any `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// APIError is an error that stores the StatusCode and Status from the response.
type APIError struct {
	StatusCode int
	Status     string
}

// Error implements the Error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("api error: %d %s", e.StatusCode, e.Status)
}

// NewAPIError creates an instance of APIError from a resty.Response. It extracts the
// StatusCode and Status from the response.
func NewAPIError(resp *resty.Response) *APIError {
	respErr := ""
	if resp.Error() != nil {
		errResponse := resp.Error().(*ErrorResponse)
		respErr = errResponse.Error
		if respErr == "" {
			respErr = string(resp.Body())
		}
	}
	return &APIError{
		StatusCode: resp.StatusCode(),
		Status:     fmt.Sprintf("%s: %s", resp.Status(), respErr),
	}
}

// Used for now to resolve the checking
var tlsConfig = tls.Config{InsecureSkipVerify: true}

// Client is REST client for the Materials Commons API.
type Client struct {
	APIKey  string
	BaseURL string
	rClient *resty.Client
}

// ClientArgs are the arguments when creating the client. You specify the URL to the server and the
// API Key for the user. If BaseURL is blank then it defaults to https://materialscommons.org/api.
type ClientArgs struct {
	APIKey  string
	BaseURL string
}

// NewClient creates a new client, sets the Accept and Content-Type headers to
// "application/json", and sets the Authorization header to the token. It
// does a small amount of cleaning on the BaseURL by removing the trailing
// slashes in the baseURL so the API can construct paths easier.
func NewClient(args *ClientArgs) *Client {
	c := resty.New().
		SetTLSClientConfig(&tlsConfig).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetAuthToken(args.APIKey)

	baseURL := "https://materialscommons.org/api"
	if args.BaseURL != "" {
		baseURL = strings.TrimSuffix(args.BaseURL, "/")
	}
	return &Client{
		BaseURL: baseURL,
		APIKey:  args.APIKey,
		rClient: c,
	}
}

func (c *Client) SetDebug(on bool) {
	c.rClient.SetDebug(on)
}

func checkError(resp *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.IsError() {
		return NewAPIError(resp)
	}
	return nil
}

func (c *Client) r() *resty.Request {
	return c.rClient.R()
}

func (c *Client) CreateProject(req CreateProjectRequest) (*mcmodel.Project, error) {
	proj := &mcmodel.Project{}

	url := c.BaseURL + "/projects"
	resp, err := c.r().
		SetBody(req).
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{proj}).
		Post(url)

	if err := checkError(resp, err); err != nil {
		return nil, err
	}

	return proj, nil
}

func (c *Client) GetProject(id int) (*mcmodel.Project, error) {
	proj := &mcmodel.Project{}

	url := c.BaseURL + fmt.Sprintf("/projects/%d", id)
	resp, err := c.r().
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{proj}).
		Get(url)

	if err := checkError(resp, err); err != nil {
		return nil, err
	}

	return proj, nil
}

func (c *Client) DeleteProject(id int) error {
	url := c.BaseURL + fmt.Sprintf("/projects/%d", id)
	resp, err := c.r().Delete(url)
	return checkError(resp, err)
}

func (c *Client) CreateExperiment(request CreateExperimentRequest) (*mcmodel.Experiment, error) {
	experiment := &mcmodel.Experiment{}

	url := c.BaseURL + "/experiments"
	resp, err := c.r().
		SetBody(request).
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{experiment}).
		Post(url)

	if err := checkError(resp, err); err != nil {
		return nil, err
	}
	return experiment, nil
}

func (c *Client) CreateDataset(projectID int, req CreateOrUpdateDatasetRequest) (*mcmodel.Dataset, error) {
	dataset := &mcmodel.Dataset{}

	url := c.BaseURL + fmt.Sprintf("/projects/%d/datasets", projectID)
	resp, err := c.r().
		SetBody(req).
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{dataset}).
		Post(url)

	if err := checkError(resp, err); err != nil {
		return nil, err
	}
	return dataset, nil
}

func (c *Client) GetDataset(projectID int, datasetID int) (*mcmodel.Dataset, error) {
	dataset := &mcmodel.Dataset{}

	url := c.BaseURL + fmt.Sprintf("/projects/%d/datasets/%d", projectID, datasetID)
	resp, err := c.r().
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{dataset}).
		Get(url)

	if err := checkError(resp, err); err != nil {
		return nil, err
	}
	return dataset, nil
}

func (c *Client) UpdateDataset(projectID int, datasetID int, req CreateOrUpdateDatasetRequest) (*mcmodel.Dataset, error) {
	dataset := &mcmodel.Dataset{}

	url := c.BaseURL + fmt.Sprintf("/projects/%d/datasets/%d", projectID, datasetID)
	resp, err := c.r().
		SetBody(req).
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{dataset}).
		Put(url)

	if err := checkError(resp, err); err != nil {
		return nil, err
	}
	return dataset, nil
}

func (c *Client) UpdateDatasetFileSelection(projectID, datasetID int, fileSelection DatasetFileSelection) (*mcmodel.Dataset, error) {
	dataset := &mcmodel.Dataset{}

	if fileSelection.ExcludeFiles == nil {
		fileSelection.ExcludeFiles = []string{}
	}

	if fileSelection.IncludeFiles == nil {
		fileSelection.IncludeFiles = []string{}
	}

	if fileSelection.ExcludeDirs == nil {
		fileSelection.ExcludeDirs = []string{}
	}

	if fileSelection.IncludeDirs == nil {
		fileSelection.IncludeDirs = []string{}
	}

	url := c.BaseURL + fmt.Sprintf("/projects/%d/datasets/%d/change_file_selection", projectID, datasetID)
	resp, err := c.r().
		SetBody(fileSelection).
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{dataset}).
		Put(url)

	if err := checkError(resp, err); err != nil {
		return nil, err
	}
	return dataset, nil
}

func (c *Client) PublishDataset(projectID int, datasetID int) (*mcmodel.Dataset, error) {
	dataset := &mcmodel.Dataset{}
	var req struct {
		ProjectID int `json:"project_id"`
	}
	req.ProjectID = projectID

	url := c.BaseURL + fmt.Sprintf("/datasets/%d/publish", datasetID)
	resp, err := c.r().
		SetBody(req).
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{dataset}).
		Put(url)

	if err := checkError(resp, err); err != nil {
		return nil, err
	}
	return dataset, nil
}

func (c *Client) UnpublishDataset(projectID int, datasetID int) (*mcmodel.Dataset, error) {
	dataset := &mcmodel.Dataset{}
	var req struct {
		ProjectID int `json:"project_id"`
	}
	req.ProjectID = projectID

	url := c.BaseURL + fmt.Sprintf("/datasets/%d/unpublish", datasetID)
	resp, err := c.r().
		SetBody(req).
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{dataset}).
		Put(url)

	if err := checkError(resp, err); err != nil {
		fmt.Printf(resp.String())
		return nil, err
	}
	return dataset, nil
}

func (c *Client) CreateActivity(req CreateActivityRequest) (*mcmodel.Activity, error) {
	activity := &mcmodel.Activity{}

	url := c.BaseURL + "/activities"
	resp, err := c.r().
		SetBody(req).
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{activity}).
		Post(url)

	if err := checkError(resp, err); err != nil {
		return nil, err
	}
	return activity, nil
}

func (c *Client) CreateEntity(req CreateEntityRequest) (*mcmodel.Entity, error) {
	entity := &mcmodel.Entity{}

	if req.Category == "" {
		req.Category = "experimental"
	}

	if req.Category != "experimental" && req.Category != "computational" {
		return nil, fmt.Errorf("category must be either 'experimental' or 'computational'")
	}

	url := c.BaseURL + "/entities"
	resp, err := c.r().
		SetBody(req).
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{entity}).
		Post(url)

	if err := checkError(resp, err); err != nil {
		return nil, err
	}
	return entity, nil
}

func (c *Client) CreateEntityState(projectID, entityID, activityID int, req CreateEntityStateRequest) (*mcmodel.Entity, error) {
	entity := &mcmodel.Entity{}

	url := c.BaseURL + fmt.Sprintf("/projects/%d/entities/%d/activities/%d/create-entity-state", projectID, entityID, activityID)
	resp, err := c.r().
		SetBody(req).
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{entity}).
		Post(url)

	if err := checkError(resp, err); err != nil {
		return nil, err
	}
	return entity, nil
}

func (c *Client) UploadFile(projectID, directoryID int, filePath string) (*mcmodel.File, error) {
	var files [1]mcmodel.File

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	fileName := filepath.Base(filePath)

	url := c.BaseURL + fmt.Sprintf("/projects/%d/files/%d/upload", projectID, directoryID)
	resp, err := c.r().
		SetFileReader("files[]", fileName, f).
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{&files}).
		Post(url)

	if err := checkError(resp, err); err != nil {
		return nil, err
	}

	return &files[0], nil
}
