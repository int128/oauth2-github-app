# oauth2-github-app [![go](https://github.com/int128/oauth2-github-app/actions/workflows/go.yaml/badge.svg)](https://github.com/int128/oauth2-github-app/actions/workflows/go.yaml) [![Go Reference](https://pkg.go.dev/badge/github.com/int128/oauth2-github-app.svg)](https://pkg.go.dev/github.com/int128/oauth2-github-app)

This is a Go package for [authenticating with a GitHub App Installation](https://docs.github.com/en/developers/apps/authenticating-with-github-apps).
It is interoperable with `golang.org/x/oauth2` package.

## Getting Started

### Prerequisite

Set up your GitHub App Installation.

1. [Create a GitHub App](https://docs.github.com/en/developers/apps/creating-a-github-app)
1. [Download a private key of the GitHub App](https://docs.github.com/en/developers/apps/authenticating-with-github-apps)
1. [Install your GitHub App on your repository or organization](https://docs.github.com/en/developers/apps/installing-github-apps)

This package requires the following inputs:

- Private Key file
- App ID (number)
- Installation ID (number)

### Create a client

To create an OAuth2 client with GitHub App installation token,

```go
func run(ctx context.Context, appID, installationID, privateKeyFilename string) error {
	privateKey, err := oauth2githubapp.LoadPrivateKey(privateKeyFilename)
	if err != nil {
		return fmt.Errorf("could not load the private key: %w", err)
	}

	// create a client
	cfg := oauth2githubapp.Config{
		PrivateKey:     privateKey,
		AppID:          appID,
		InstallationID: installationID,
	}
	client := oauth2.NewClient(ctx, cfg.TokenSource(ctx))
}
```

For [`github.com/google/go-github`](https://github.com/google/go-github) package,

```go
	cfg := oauth2githubapp.Config{
		PrivateKey:     privateKey,
		AppID:          appID,
		InstallationID: installationID,
	}
	client := github.NewClient(oauth2.NewClient(ctx, cfg.TokenSource(ctx)))
```

For [`github.com/shurcooL/githubv4`](https://github.com/shurcooL/githubv4) package,

```go
	cfg := oauth2githubapp.Config{
		PrivateKey:     privateKey,
		AppID:          appID,
		InstallationID: installationID,
	}
	client := githubv4.NewClient(oauth2.NewClient(ctx, cfg.TokenSource(ctx)))
```

See also [example](example/main.go).

## GitHub Enterprise

To set your GitHub Enterprise server,

```go
	cfg := oauth2githubapp.Config{
		PrivateKey:     privateKey,
		AppID:          appID,
		InstallationID: installationID,
		BaseURL:        "https://api.github.example.com",
	}
	client := oauth2.NewClient(ctx, cfg.TokenSource(ctx))
```

## Contribution

This is an open source software. Feel free to contribute to it.
