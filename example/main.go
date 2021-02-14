package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/int128/oauth2-github-app/app"
	"golang.org/x/oauth2"
)

func run(ctx context.Context, appID, installationID, privateKeyName string) error {
	// create an oauth2 client
	privateKey, err := app.LoadPrivateKey(privateKeyName)
	if err != nil {
		return fmt.Errorf("could not load the private key: %w", err)
	}
	cfg := app.Config{
		PrivateKey:     privateKey,
		AppID:          appID,
		InstallationID: installationID,
	}
	client := oauth2.NewClient(ctx, cfg.TokenSource(ctx))

	// call an endpoint
	resp, err := client.Get("https://api.github.com/rate_limit")
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body error: %s", err)
	}
	log.Printf("response: status %d body %s", resp.StatusCode, string(b))
	return nil
}

func main() {
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
