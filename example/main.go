package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	oauth2githubapp "github.com/int128/oauth2-github-app"
	"golang.org/x/oauth2"
)

func run(ctx context.Context, appID, installationID, privateKeyName string) error {
	// create an oauth2 client
	privateKey, err := oauth2githubapp.LoadPrivateKey(privateKeyName)
	if err != nil {
		return fmt.Errorf("could not load the private key: %w", err)
	}
	cfg := oauth2githubapp.Config{
		PrivateKey:     privateKey,
		AppID:          appID,
		InstallationID: installationID,
	}
	client := oauth2.NewClient(ctx, cfg.TokenSource(ctx))

	// prepare a query
	var q struct {
		Query string `json:"query"`
	}
	q.Query = `
		query {
			viewer {
				login
			}
		}`
	reqBody, err := json.Marshal(&q)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	// call an endpoint
	resp, err := client.Post("https://api.github.com/graphql", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	log.Println(resp.Status)
	for k, vs := range resp.Header {
		for _, v := range vs {
			log.Printf("%s: %s", k, v)
		}
	}
	log.Println()
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		return fmt.Errorf("read body error: %s", err)
	}
	return nil
}

func main() {
	log.SetFlags(0)
	ctx := context.Background()

	appID := os.Getenv("GITHUB_APP_ID")
	installationID := os.Getenv("GITHUB_APP_INSTALLATION_ID")
	privateKeyName := os.Getenv("GITHUB_APP_PRIVATE_KEY_NAME")
	if appID == "" || installationID == "" || privateKeyName == "" {
		log.Fatalf("you need to set GITHUB_APP_ID, GITHUB_APP_INSTALLATION_ID and GITHUB_APP_PRIVATE_KEY_NAME")
	}
	if err := run(ctx, appID, installationID, privateKeyName); err != nil {
		log.Fatalf("error: %s", err)
	}
}
