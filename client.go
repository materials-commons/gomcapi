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

// CreateProject creates a new project with the specified parameters in CreateProjectRequest and returns the created project.
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

// GetProject retrieves the project details for the given project ID.
// It returns a pointer to the Project object and an error, if any.
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

// DeleteProject deletes a project identified by the provided ID.
// It makes a DELETE request to the API endpoint corresponding to the given project ID.
// Returns an error if the project could not be deleted or if there is any issue with the request.
func (c *Client) DeleteProject(id int) error {
	url := c.BaseURL + fmt.Sprintf("/projects/%d", id)
	resp, err := c.r().Delete(url)
	return checkError(resp, err)
}

// CreateExperiment creates a new experiment based on the given CreateExperimentRequest.
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

// CreateDataset creates a new dataset within the specified project.
// It takes a projectID and a CreateOrUpdateDatasetRequest as parameters.
// It returns a pointer to the created Dataset object or an error, if any occurs.
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

// UpdateDataset updates an existing dataset for the given project.
// Takes in a projectID, datasetID, and a CreateOrUpdateDatasetRequest object.
// Returns the updated Dataset object or an error if the update fails.
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

// UpdateDatasetFileSelection updates the file selection criteria for a specified dataset within a project.
// It includes and excludes specified files and directories based on the DatasetFileSelection object.
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

// PublishDataset publishes a specified dataset by its datasetID in a particular project identified by projectID.
// Returns the published Dataset object or an error if the operation fails.
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

// UnpublishDataset unpublishes a dataset associated with a specified project and dataset ID,
// returning the updated dataset or an error if the operation fails.
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
		return nil, err
	}
	return dataset, nil
}

// CreateActivity creates a new activity based on the provided CreateActivityRequest struct.
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

// CreateEntity creates a new entity based on the provided request and returns the created entity or an error.
// The category in the request must be either 'experimental' or 'computational'. Defaults to 'experimental'.
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

// CreateEntityState creates a new entity state associated with the provided project, entity, and activity IDs.
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

// GetFileByPath fetches a file from a project given the specified project ID and file path.
// It returns a pointer to the file and an error if the fetch operation fails.
func (c *Client) GetFileByPath(projectID int, path string) (*mcmodel.File, error) {
	file := &mcmodel.File{}
	req := struct {
		ProjectID int    `json:"project_id"`
		Path      string `json:"path"`
	}{
		ProjectID: projectID,
		Path:      path,
	}

	url := c.BaseURL + "/files/by_path"
	resp, err := c.r().
		SetBody(req).
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{file}).
		Post(url)
	if err := checkError(resp, err); err != nil {
		return nil, err
	}
	return file, nil
}

// CreateDirectoryByPath creates a directory at the specified path within the given project. If the
// directory already exists, it returns the existing directory. It takes a project ID and a path as
// parameters and returns the created directory or an error.
func (c *Client) CreateDirectoryByPath(projectID int, path string) (*mcmodel.File, error) {
	file := &mcmodel.File{}
	req := struct {
		ProjectID int    `json:"project_id"`
		Path      string `json:"path"`
	}{
		ProjectID: projectID,
		Path:      path,
	}

	resp, err := c.r().
		SetBody(req).
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{file}).
		Post(c.BaseURL + "/directories/by-path")
	if err := checkError(resp, err); err != nil {
		return nil, err
	}
	return file, nil
}

// UploadFileTo uploads a file to a specified project directory path. If the project directory path is
// blank then it uploads the file to the project root. If the directory does not exist, it creates it.
func (c *Client) UploadFileTo(projectID int, filePath string, projectPath string) (*mcmodel.File, error) {
	if projectPath == "" {
		projectPath = "/"
	}
	dir, err := c.CreateDirectoryByPath(projectID, projectPath)
	if err != nil {
		return nil, err
	}

	return c.UploadFile(projectID, dir.ID, filePath)
}

// UploadFile uploads a file to the specified project and directory.
// Parameters:
// - projectID: ID of the project to which the file will be uploaded.
// - directoryID: ID of the directory within the project where the file will be stored.
// - filePath: Local path of the file to be uploaded.
// Returns:
// - *mcmodel.File: Uploaded file metadata if the operation is successful.
// - error: Describes the error encountered during file upload.
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

// DepositDataset deposits a dataset in a specified project given the project ID and request details.
// It creates the dataset, uploads files into a unique directory, and sets file selection for the dataset.
func (c *Client) DepositDataset(projectID int, req DepositDatasetRequest) (*mcmodel.Dataset, error) {
	// 1. Create the dataset.
	createDatasetReq := CreateOrUpdateDatasetRequest{
		Name:        req.Metadata.Name,
		Description: req.Metadata.Description,
		Summary:     req.Metadata.Summary,
		License:     req.Metadata.License,
		Funding:     req.Metadata.Funding,
		//Experiments: nil,
		//Communities: nil,
		Tags:    req.Metadata.Tags,
		Authors: req.Metadata.Authors,
	}
	dataset, err := c.CreateDataset(projectID, createDatasetReq)
	if err != nil {
		return nil, err
	}

	// 2. Add the additional metadata to the dataset
	// 3. Upload the files
	//    Create a unique directory for the dataset in the project?
	//    Or should we instead create a project per dataset?

	// For now lets create a unique directory. This makes dataset file selection easy.
	// The directory will have the dataset UUID as its name. This is kind of a crappy
	// solution, but for now let's start with that. We can revise this decision
	// later on after discussing with Valentin.
	dir, err := c.CreateDirectoryByPath(projectID, "/"+dataset.UUID)
	if err != nil {
		return nil, err
	}

	// Now for each of the files we are uploading, we need to append
	// the created directory to its path.
	for _, file := range req.Files {
		fileDir := file.Directory
		if fileDir == "" {
			fileDir = "/"
		} else {
			fileDir = fileDir + "/"
		}
		_, err := c.UploadFileTo(projectID, file.File, dir.Path+file.Directory)
		if err != nil {
			// For now lets stop all uploads and return an error
			return nil, err
		}
	}

	// 4. Set the dataset file selection
	fileSelection := DatasetFileSelection{
		IncludeFiles: nil,
		ExcludeFiles: nil,
		IncludeDirs:  []string{"/" + dataset.UUID},
		ExcludeDirs:  nil,
	}
	dataset, err = c.UpdateDatasetFileSelection(projectID, dataset.ID, fileSelection)
	if err != nil {
		return nil, err
	}

	return dataset, nil
}

// ListProjects lists all the projects a user is a member of
func (c *Client) ListProjects() ([]mcmodel.Project, error) {
	var projects []mcmodel.Project
	url := c.BaseURL + "/projects"
	resp, err := c.r().
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{&projects}).
		Get(url)
	if err := checkError(resp, err); err != nil {
		return nil, err
	}
	return projects, nil
}

// ListDatasets lists all the datasets in a project
func (c *Client) ListDatasets(projectID int) ([]mcmodel.Dataset, error) {
	var datasets []mcmodel.Dataset
	url := c.BaseURL + fmt.Sprintf("/projects/%d/datasets", projectID)
	resp, err := c.r().
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{&datasets}).
		Get(url)
	if err := checkError(resp, err); err != nil {
		return nil, err
	}
	return datasets, nil
}

// MintDOIForDataset mints a new (findable) DOI for the dataset and assigns the DOI to it.
func (c *Client) MintDOIForDataset(projectID, datasetID int) (*mcmodel.Dataset, error) {
	var dataset mcmodel.Dataset
	url := c.BaseURL + fmt.Sprintf("/projects/%d/datasets/%d/assign_doi", projectID, datasetID)
	resp, err := c.r().
		SetError(&ErrorResponse{}).
		SetResult(&DataWrapper{&dataset}).
		Post(url)
	if err := checkError(resp, err); err != nil {
		return nil, err
	}
	return &dataset, nil
}
