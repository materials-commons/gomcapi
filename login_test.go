package mcapi

import (
	"testing"

	"github.com/materials-commons/gomcapi/pkg/tutils/assert"
)

func TestLogin(t *testing.T) {
	c, err := Login("test@test.mc", "test", "http://mcdev.localhost/api")
	assert.Okf(t, err, "Login failed with err :%s", err)
	assert.NotNil(t, c)
	assert.Equals(t, c.APIKey, "totally-bogus")
}
