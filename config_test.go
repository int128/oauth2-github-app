package oauth2githubapp

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// https://docs.github.com/en/rest/apps/apps?apiVersion=2022-11-28#create-an-installation-access-token-for-an-app
const fakeAccessToken = `{
  "token": "ghs_16C7e42F292c6912E7710c838347Ae178B4a",
  "expires_at": "2016-07-11T22:14:10Z",
  "permissions": {
    "issues": "write",
    "contents": "read"
  },
  "repository_selection": "selected",
  "repositories": []
}`

type fakeCreateInstallationAccessToken struct {
	t         *testing.T
	publicKey *rsa.PublicKey
	appID     string
}

func (f fakeCreateInstallationAccessToken) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		f.t.Errorf("method wants POST but was %s", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	authorizationHeader := r.Header.Get("authorization")
	if err := f.authorize(authorizationHeader); err != nil {
		f.t.Logf("authorization header: %s", authorizationHeader)
		f.t.Errorf("invalid authorization header: %s", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(201)
	if _, err := io.WriteString(w, fakeAccessToken); err != nil {
		f.t.Errorf("write error: %s", err)
	}
}

func (f fakeCreateInstallationAccessToken) authorize(header string) error {
	if !strings.HasPrefix(header, "Bearer ") {
		return fmt.Errorf("header wants Bearer token")
	}
	token := strings.TrimPrefix(header, "Bearer ")
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return fmt.Errorf("token wants 3 parts but was %d parts", len(parts))
	}
	hashed := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	sign, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return fmt.Errorf("invalid base64 of sign: %w", err)
	}
	if err := rsa.VerifyPKCS1v15(f.publicKey, crypto.SHA256, hashed[:], sign); err != nil {
		return fmt.Errorf("verify error: %s", err)
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("invalid base64 of payload: %w", err)
	}
	var payloadToken struct {
		Issuer string `json:"iss"`
	}
	if err := json.Unmarshal(payload, &payloadToken); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}
	if f.appID != payloadToken.Issuer {
		return fmt.Errorf("issuer wants %s but was %s", f.appID, payloadToken.Issuer)
	}
	return nil
}

func TestConfig_Token(t *testing.T) {
	privateKey, err := LoadPrivateKey("testdata/oauth2-github-app-test.2021-02-12.private-key.pem")
	if err != nil {
		t.Fatalf("could not load the private key: %s", err)
	}

	fakeGitHub := http.NewServeMux()
	fakeGitHub.Handle("/app/installations/MyInstallationID/access_tokens", fakeCreateInstallationAccessToken{
		t:         t,
		publicKey: &privateKey.PublicKey,
		appID:     "MyAppID",
	})
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
		t.Errorf("token.AccessToken wants %s but was %s", w, token.AccessToken)
	}
	if w := time.Date(2016, 7, 11, 22, 14, 10, 0, time.UTC); token.Expiry != w {
		t.Errorf("token.Expiry wants %s but was %s", w, token.Expiry)
	}
	if w := "token"; token.TokenType != w {
		t.Errorf("token.TokenType wants %s but was %s", w, token.AccessToken)
	}
}
