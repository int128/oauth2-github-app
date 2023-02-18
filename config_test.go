package oauth2githubapp

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeGitHubAppInstallationTokenAPI struct{}

func (f fakeGitHubAppInstallationTokenAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method should be POST", http.StatusMethodNotAllowed)
		return
	}
	authorizationHeader := r.Header.Get("authorization")
	if authorizationHeader == "" {
		http.Error(w, "No authorization header", http.StatusUnauthorized)
		return
	}
	// https://docs.github.com/en/rest/apps/apps?apiVersion=2022-11-28#create-an-installation-access-token-for-an-app
	w.WriteHeader(201)
	_, _ = io.WriteString(w, `{
  "token": "ghs_16C7e42F292c6912E7710c838347Ae178B4a",
  "expires_at": "2016-07-11T22:14:10Z",
  "permissions": {
    "issues": "write",
    "contents": "read"
  },
  "repository_selection": "selected",
  "repositories": []
}`)
}

func TestConfig_Token(t *testing.T) {
	privateKey, err := LoadPrivateKey("testdata/oauth2-github-app-test.2021-02-12.private-key.pem")
	if err != nil {
		t.Fatalf("could not load the private key: %s", err)
	}

	fakeGitHub := http.NewServeMux()
	fakeGitHub.Handle("/app/installations/MyInstallationID/access_tokens", fakeGitHubAppInstallationTokenAPI{})
	fakeServer := httptest.NewServer(fakeGitHub)
	defer fakeServer.Close()

	cfg := Config{
		AppID:          "MyAppID",
		InstallationID: "MyInstallationID",
		PrivateKey:     privateKey,
		BaseURL:        fakeServer.URL,
	}
	token, err := cfg.Token(context.TODO())
	if err != nil {
		t.Fatalf("token error: %s", err)
	}
	if w := "ghs_16C7e42F292c6912E7710c838347Ae178B4a"; token.AccessToken != w {
		t.Errorf("token wants %s but was %s", w, token.AccessToken)
	}
}
