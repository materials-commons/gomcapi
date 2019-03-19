package mcapi

import (
	"github.com/materials-commons/config"
	"gopkg.in/resty.v1"
)

// projectsRoute creates the default route path for projects REST API
func projectsRoute() string {
	p := config.GetString("mcurl") + "/v2" + "/projects"
	return p
}

// r is similar to resty.R() except that it sets the TLS configuration and the apikey
func r() *resty.Request {
	return resty.SetTLSClientConfig(&tlsConfig).R().SetQueryParam("apikey", config.GetString("apikey"))
}

// GetAllProjects gets all projects that the user has access to
func GetAllProjects() ([]*Project, error) {
	var results []*Project
	_, err := r().SetResult(&results).Get(projectsRoute())
	return results, err
}

// GetProject retrieves the given from with projectID
func GetProject(projectID string) (*Project, error) {
	var proj Project

	_, err := r().SetResult(&proj).Get(projectsRoute() + "/" + projectID)
	return &proj, err
}

// CreateProject creates a new project with the user set as the owner
func CreateProject(name, description string) (*Project, error) {
	var proj Project
	_, err := r().SetResult(&proj).
		SetBody(map[string]interface{}{"name": name, "description": description}).
		Post(projectsRoute())
	return &proj, err
}

func (p *Project) Update() error {
	return nil
}

func (p *Project) Delete() error {
	return nil
}

func (p *Project) CreateExperiment(name, description string) (*Experiment, error) {
	var experiment Experiment
	route := projectsRoute() + "/" + p.ID + "/experiments"
	_, err := r().SetResult(&experiment).SetBody(map[string]interface{}{
		"name":        name,
		"description": description,
		"aim":         "",
		"status":      "active",
		"tasks":       []string{},
	}).
		Post(route)
	return &experiment, err
}

func (e *Experiment) DeleteExperiment() error {
	return nil
}

// DeleteProject deletes the projects with the given projectID
func DeleteProject(projectID string) error {
	_, err := r().Delete(projectsRoute() + "/" + projectID)
	return err
}