package mcapi

import (
	"crypto/tls"
	"fmt"

	"github.com/materials-commons/config"
	"gopkg.in/resty.v1"
)

type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Owner     string    `json:"owner"`
	Birthtime Timestamp `json:"birthtime"`
	MTime     Timestamp `json:"mtime"`
}

func projectsRoute() string {
	fmt.Println(config.GetString("mcurl"))
	p := config.GetString("mcurl") + "/v2" + "/projects"
	return p
}

var tlsConfig = tls.Config{InsecureSkipVerify: true}

func p() *resty.Request {
	return resty.SetTLSClientConfig(&tlsConfig).R().SetQueryParam("apikey", config.GetString("apikey"))
}

func GetAllProjects() ([]*Project, error) {
	var results []*Project
	_, err := p().SetResult(&results).Get(projectsRoute())
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (p *Project) Update() error {
	return nil
}

func (p *Project) Delete() error {
	return nil
}
