package dreamhost

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	dreamhostLiveTest bool
	dreamhostApiKey   string
	testDomain        string
)

func init() {
	dreamhostApiKey = os.Getenv("DREAMHOST_API_KEY")
	testDomain = os.Getenv("DREAMHOST_DOMAIN")

	if len(dreamhostApiKey) > 0 {
		dreamhostLiveTest = true
	}
}

func restoreDreamhostEnv() {
	os.Setenv("DREAMHOST_API_KEY", dreamhostApiKey)
}

func TestNewDNSProviderValidEnv(t *testing.T) {
	if !dreamhostLiveTest {
		t.Skip("skipping live test (requires credentials)")
	}

	os.Setenv("DREAMHOST_API_KEY", "other")
	_, err := NewDNSProvider()
	assert.NoError(t, err)
	restoreDreamhostEnv()
}

func TestNewDNSProviderMissingCredErr(t *testing.T) {
	os.Setenv("DREAMHOST_API_KEY", "")
	_, err := NewDNSProvider()
	assert.EqualError(t, err, "Dreamhost credentials missing")
	restoreDreamhostEnv()
}

func TestLiveDreamhostPresentCleanUp(t *testing.T) {
	if !dreamhostLiveTest {
		t.Skip("skipping live test (requires credentials)")
	}

	provider, err := NewDNSProvider()
	assert.NoError(t, err)

	err = provider.Present(testDomain, "", "123d==")
	assert.NoError(t, err)

	err = provider.CleanUp(testDomain, "", "123d==")
	assert.NoError(t, err)
}
