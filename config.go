// Package oauth2githubapp provides an OAuth 2.0 client for GitHub App.
//
// This package implements the authentication method described in
// https://docs.github.com/en/developers/apps/authenticating-with-github-apps
package oauth2githubapp

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

const defaultBaseURL = "https://api.github.com"

// Config represents a config of GitHub App Installation
type Config struct {
	PrivateKey     *rsa.PrivateKey
	AppID          string
	InstallationID string

	// BaseURL is an endpoint of GitHub API.
	// Any trailing slash is trimmed.
	// If not set, it defaults to https://api.github.com
	BaseURL string
}

// TokenSource returns an oauth2.TokenSource for GitHub App Installation
func (c Config) TokenSource(ctx context.Context) oauth2.TokenSource {
	return oauth2.ReuseTokenSource(nil, &tokenRefresher{ctx: ctx, config: c})
}

type tokenRefresher struct {
	ctx    context.Context
	config Config
}

func (r tokenRefresher) Token() (*oauth2.Token, error) {
	token, err := r.config.Token(r.ctx)
	if err != nil {
		return nil, fmt.Errorf("could not refresh an access token: %w", err)
	}
	return token, nil
}

// Token requests an installation access token using a JWT of GitHub App.
//
// See https://docs.github.com/en/developers/apps/authenticating-with-github-apps#authenticating-as-a-github-app
func (c Config) Token(ctx context.Context) (*oauth2.Token, error) {
	appJWT, err := newAppJWT(c.AppID, c.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("could not generate a jwt: %w", err)
	}
	resp, err := c.createInstallationAccessToken(ctx, appJWT)
	if err != nil {
		return nil, fmt.Errorf("could not create an installation access token: %w", err)
	}
	return &oauth2.Token{
		AccessToken: resp.Token,
		Expiry:      resp.ExpiresAt,
		TokenType:   "token",
	}, nil
}

type tokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// https://docs.github.com/en/rest/apps/apps?apiVersion=2022-11-28#create-an-installation-access-token-for-an-app
func (c Config) createInstallationAccessToken(ctx context.Context, appJWT string) (*tokenResponse, error) {
	baseURL := strings.TrimRight(c.BaseURL, "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	tokenEndpoint := fmt.Sprintf("%s/app/installations/%s/access_tokens", baseURL, url.PathEscape(c.InstallationID))
	req, err := http.NewRequestWithContext(ctx, "POST", tokenEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", appJWT))
	req.Header.Add("accept", "application/vnd.github.v3+json")

	hc, ok := ctx.Value(oauth2.HTTPClient).(*http.Client)
	if !ok {
		hc = http.DefaultClient
	}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("http status %d, body error: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("http status %d, body %s", resp.StatusCode, string(b))
	}

	d := json.NewDecoder(resp.Body)
	var t tokenResponse
	if err := d.Decode(&t); err != nil {
		return nil, fmt.Errorf("invalid json response: %w", err)
	}
	return &t, nil
}
