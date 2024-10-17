package mcapi

import (
	"fmt"
	"testing"
)

func TestDeposit_ExampleDepositDataset(t *testing.T) {
	c := newClientForTest()

	// First we will create a project to store our workflow in.
	// Since projects are scheduled for deletion we will add a
	// random string to the name.

	r := generateRandomString(5)
	req := CreateProjectRequest{
		Name:        "example_deposit_" + r,
		Description: "test project",
		Summary:     "test project",
	}

	proj, err := c.CreateProject(req)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", proj)

	depositReq := DepositDatasetRequest{
		Files: []DatasetFileUpload{
			{
				Description: "An example file",
				File:        "/tmp/example.txt",
				Directory:   "",
			},
		},
		Metadata: DatasetMetadata{
			Name:        "Example Dataset Deposit",
			Description: "Example of depositing a dataset",
			Summary:     "Example of depositing a dataset",
			License:     DatasetLicensePublicDomain,
			Funding:     "DOE",
			Communities: nil,
			Authors:     nil,
			Tags:        []Tag{{Value: "example"}},
			DOI:         "",
			Papers:      []Paper{{Name: "Example", DOI: "", Reference: ""}},
			Attributes:  nil,
		},
	}

	dataset, err := c.DepositDataset(proj.ID, depositReq)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", dataset)
}
