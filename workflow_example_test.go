package mcapi

import (
	"fmt"
	"testing"
	"time"
)

func TestWorkflow_ExampleCreateEntityState(t *testing.T) {
	c := newClientForTest()

	// First we will create a project to store our workflow in.
	// Since projects are scheduled for deletion we will add a
	// random string to the name.

	r := generateRandomString(5)
	req := CreateProjectRequest{
		Name:        "example_workflow" + r,
		Description: "test project",
		Summary:     "test project",
	}

	proj, err := c.CreateProject(req)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", proj)

	f, err := c.UploadFile(proj.ID, proj.RootDir.ID, "/home/gtarcea/Downloads/upload.txt")

	if err != nil {
		_ = c.DeleteProject(proj.ID)
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", f)

	// Create dataset
	dsReq := CreateOrUpdateDatasetRequest{
		Name:        "ds1",
		Description: "ds1 dataset",
	}

	ds, err := c.CreateDataset(proj.ID, dsReq)
	if err != nil {
		_ = c.DeleteProject(proj.ID)
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", ds)

	// Add files to dataset
	fSelection := DatasetFileSelection{
		IncludeFiles: []string{
			"/upload.txt",
		},
	}

	ds, err = c.UpdateDatasetFileSelection(proj.ID, ds.ID, fSelection)
	if err != nil {
		_ = c.DeleteProject(proj.ID)
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", ds)

	ds, err = c.PublishDataset(proj.ID, ds.ID)

	if err != nil {
		_ = c.DeleteProject(proj.ID)
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", ds)

	// Give time to complete publishing the dataset
	for {
		time.Sleep(5 * time.Second)
		ds, err = c.GetDataset(proj.ID, ds.ID)
		if !ds.PublishedAt.IsZero() {
			break
		}
	}

	// Clean up by removing project

	// Step 1 - Can't delete a project that has a published dataset.
	ds, err = c.UnpublishDataset(proj.ID, ds.ID)

	if err != nil {
		t.Fatal(err)
	}

	// Give time to complete unpublishing the dataset
	for {
		time.Sleep(5 * time.Second)
		ds, err = c.GetDataset(proj.ID, ds.ID)
		if ds.PublishedAt.IsZero() {
			break
		}
	}

	// Step 2 - Delete the project
	if err := c.DeleteProject(proj.ID); err != nil {
		t.Error(err)
	}
}
