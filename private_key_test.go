package oauth2githubapp_test

import (
	"testing"

	"github.com/int128/oauth2-github-app"
)

func TestLoadPrivateKey(t *testing.T) {
	privateKey, err := oauth2githubapp.LoadPrivateKey("testdata/oauth2-github-app-test.2021-02-12.private-key.pem")
	if err != nil {
		t.Fatalf("load error: %s", err)
	}
	size := privateKey.Size()
	if size != 256 {
		t.Errorf("size wants 256 but was %d", size)
	}
}
