# oauth2-github-app [![Go Reference](https://pkg.go.dev/badge/github.com/int128/oauth2-github-app.svg)](https://pkg.go.dev/github.com/int128/oauth2-github-app)

This is a Go package for [authenticating with a GitHub App Installation](https://docs.github.com/en/developers/apps/authenticating-with-github-apps).

Features of this package:

- Interoperable with `golang.org/x/oauth2` package
- Least dependency

Here are advantages of the installation access token of GitHub App:

- Personal access token
  - Belong to a user
  - Share the rate limit in a user
- Installation access token
  - Belong to a user or organization
  - Have each rate limit for an installation


## Getting Started

You need to set up your GitHub App Installation as follows:

1. [Create a GitHub App](https://docs.github.com/en/developers/apps/creating-a-github-app)
1. [Download a private key of the GitHub App](https://docs.github.com/en/developers/apps/authenticating-with-github-apps)
1. [Install your GitHub App on your repository or organization](https://docs.github.com/en/developers/apps/installing-github-apps)

Finally you will get the following items:

- App ID (integer)
- Private Key (file)
- Installation ID (integer)

To create a http client with authentication:

```go
package main

import (
	"context"
	"fmt"

	"github.com/int128/oauth2-github-app"
	"golang.org/x/oauth2"
)

func run(ctx context.Context, appID, installationID, privateKeyName string) error {
	// create an http client
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
	// ...
}
```

See also the [example application](example/main.go).


## Contribution

This is an open source software. Feel free to contribute to it.
