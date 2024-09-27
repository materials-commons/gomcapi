package mcapi

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func newClientForTest() *Client {
	args := &ClientArgs{
		BaseURL: "http://localhost:8000/api",
		APIKey:  os.Getenv("MCAPI_KEY"),
	}

	return NewClient(args)
}

func TestClient_CreateProject(t *testing.T) {
	c := newClientForTest()

	req := CreateProjectRequest{
		Name:        generateRandomString(10),
		Description: "test project",
		Summary:     "test project",
	}

	proj, err := c.CreateProject(req)

	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", proj)
}

func TestClient_GetProject(t *testing.T) {
	c := newClientForTest()
	proj, err := c.GetProject(438)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", proj)
}

func TestClient_CreateDataset(t *testing.T) {
	c := newClientForTest()
	req := CreateOrUpdateDatasetRequest{
		Name:        generateRandomString(10),
		Description: "test dataset",
		Summary:     "test dataset",
	}
	ds, err := c.CreateDataset(438, req)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", ds)
}

func TestClient_CreateActivity(t *testing.T) {
	c := newClientForTest()
	req := CreateActivityRequest{
		Name:         generateRandomString(10),
		Description:  "test activity",
		ProjectID:    438,
		ExperimentID: 879,
		Attributes: []Attribute{
			{
				Name:  "attr1",
				Value: 1,
				Unit:  "c",
			},
		},
	}

	activity, err := c.CreateActivity(req)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", activity)
}

func TestClient_CreateEntity(t *testing.T) {
	c := newClientForTest()
	req := CreateEntityRequest{
		Name:         generateRandomString(10),
		Description:  "test entity",
		ProjectID:    438,
		ExperimentID: 879,
		ActivityID:   230380,
		Attributes: []Attribute{
			{
				Name:  "attr1",
				Value: 1,
				Unit:  "c",
			},
		},
	}
	entity, err := c.CreateEntity(req)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", entity)
}

func TestClient_CreateEntityState(t *testing.T) {
	c := newClientForTest()
	req := CreateActivityRequest{
		Name:         generateRandomString(10),
		Description:  "step 2",
		ProjectID:    438,
		ExperimentID: 879,
		Attributes: []Attribute{
			{
				Name:  "attr1",
				Value: 1,
				Unit:  "c",
			},
		},
	}

	activity, err := c.CreateActivity(req)
	if err != nil {
		t.Error(err)
	}

	req2 := CreateEntityStateRequest{
		Current: true,
		Attributes: []Attribute{
			{
				Name:  "attr2",
				Value: 1,
				Unit:  "m",
			},
		},
	}

	entity, err := c.CreateEntityState(438, 50115, activity.ID, req2)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", entity)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
