package authenticate

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
)

func Authenticate(access_token string) (context.Context, *github.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: access_token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return ctx, client
}
