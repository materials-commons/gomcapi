package mcapi_test

import (
	"testing"

	"github.com/materials-commons/gomcapi/pkg/tutils/assert"

	. "github.com/materials-commons/gomcapi"
)

func TestGetAllProjects(t *testing.T) {
	projects, err := GetAllProjects()
	assert.Ok(t, err)
	assert.NotNil(t, projects)
}

func TestCreateProjectAndExperiments(t *testing.T) {
	var (
		proj *Project
		err  error
		e    *Experiment
	)
	projName := "Proj1"
	projDescription := "Project Created With Test"
	proj, err = CreateProject(projName, projDescription)
	assert.Ok(t, err)
	assert.NotNil(t, proj)
	assert.Equals(t, projName, proj.Name)
	assert.Equals(t, projDescription, proj.Description)

	e, err = proj.CreateExperiment("t1", "t1 description")
	assert.Ok(t, err)
	assert.NotNil(t, e)

	err = DeleteProject(proj.ID)
	assert.Ok(t, err)
}
