package main

import (
	"os"

	"golang.org/x/net/context"
	"github.com/google/go-github/github"
)

func SpoofUser(last_user *github.User) {
	ExecCommand("git", "config", "user.email", *last_user.Email)
	ExecCommand("git", "config", "user.name", *last_user.Name)
}

func Clear() {
	ExecCommand("git", "cherry-pick", "--abort")
	ExecCommand("git", "rebase", "--abort")
}

func Fetch(PR *github.PullRequest) {
	creator := *PR.User.Login
	ExecCommand("git", "remote", "add", creator, *PR.Head.Repo.GitURL)
	ExecCommand("git", "fetch", creator)
}

func CheckoutBranch(branch string) error {
	if ExecCommand("git", "checkout", "-b", branch) != nil {
		return ExecCommand("git", "checkout", branch)
	} else {
		return nil
	}
}

func CherryPick(PR *github.PullRequest) error {
	return ExecCommand("git", "cherry-pick", *PR.Base.SHA + ".." + *PR.Head.SHA)
}

func Push(login string, project string, branch string) {
	ExecCommand("git", "push", "--set-upstream", "https://github.com/" + login + "/" + project, branch, "--force")
}

func OpenPR(client *github.Client, ctx context.Context, login string, project string, head string) *github.PullRequest {
	title := "cherry-pick-bot with a bunch of commits"
	base := "master"

	open_PR, _, err := client.PullRequests.Create(
		ctx, login, project,
		&github.NewPullRequest{Title: &title, Head: &head, Base: &base})
	Die(err)

	return open_PR
}

func Rebase(branch string) error {
	return ExecCommand("git", "pull", "--rebase", "origin", branch)
}

func ChangeRepo(login string, project string) {
	os.MkdirAll(WorkDir + login, 0775)
	os.Chdir(WorkDir + login)
	ExecCommand("git", "clone", "git://github.com/" + login + "/" + project)
	os.Chdir(WorkDir + login + "/" + project)
}
