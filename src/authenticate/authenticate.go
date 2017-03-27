package authenticate

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"

	"config"
)

func Authenticate() (context.Context, *github.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return ctx, client
}
