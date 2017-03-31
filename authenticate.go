package main

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
)

func authenticate() (context.Context, *github.Client) {
	logger.Info("Setting up Github client ...")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return ctx, client
}
