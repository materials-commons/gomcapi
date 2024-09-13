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

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
