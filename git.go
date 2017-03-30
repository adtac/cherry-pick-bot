package main

import (
	"os"

	"golang.org/x/net/context"
	"github.com/google/go-github/github"
)

func spoofUser(lastUser *github.User) {
	execCommand("git", "config", "user.email", *lastUser.Email)
	execCommand("git", "config", "user.name", *lastUser.Name)
}

func clear() {
	execCommand("git", "cherry-pick", "--abort")
	execCommand("git", "rebase", "--abort")
}

func fetch(pr *github.PullRequest) {
	creator := *pr.User.Login
	execCommand("git", "remote", "add", creator, *pr.Head.Repo.GitURL)
	execCommand("git", "fetch", creator)
}

func checkoutBranch(branch string) error {
	if execCommand("git", "checkout", "-b", branch) != nil {
		return execCommand("git", "checkout", branch)
	} else {
		return nil
	}
}

func cherryPick(pr *github.PullRequest) error {
	return execCommand("git", "cherry-pick", *pr.Base.SHA + ".." + *pr.Head.SHA)
}

func push(login string, project string, branch string) {
	execCommand("git", "push", "--set-upstream", "https://github.com/" + login + "/" + project, branch, "--force")
}

func openPR(client *github.Client, ctx context.Context, login string, project string, head string) (*github.PullRequest, error) {
	title := "cherry-pick-bot with a bunch of commits"
	base := "master"

	openPR, _, err := client.PullRequests.Create(
		ctx, login, project,
		&github.NewPullRequest{Title: &title, Head: &head, Base: &base})
	if err != nil {
		return nil, err
	}


	return openPR, nil
}

func rebase(branch string) error {
	return execCommand("git", "pull", "--rebase", "origin", branch)
}

func changeRepo(login string, project string) {
	os.MkdirAll(workDir + login, 0775)
	os.Chdir(workDir + login)
	execCommand("git", "clone", "git://github.com/" + login + "/" + project)
	os.Chdir(workDir + login + "/" + project)
}
