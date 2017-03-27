package git

import (
	"os"

	"golang.org/x/net/context"
	"github.com/google/go-github/github"

	"utils"
	"config"
)

func SpoofUser(last_user *github.User) {
	utils.ExecCommand("git", "config", "user.email", *last_user.Email)
	utils.ExecCommand("git", "config", "user.name", *last_user.Name)
}

func Clear() {
	utils.ExecCommand("git", "cherry-pick", "--abort")
	utils.ExecCommand("git", "rebase", "--abort")
}

func Fetch(PR *github.PullRequest) {
	creator := *PR.User.Login
	utils.ExecCommand("git", "remote", "add", creator, *PR.Head.Repo.GitURL)
	utils.ExecCommand("git", "fetch", creator)
}

func CheckoutBranch(branch string) error {
	if utils.ExecCommand("git", "checkout", "-b", branch) != nil {
		return utils.ExecCommand("git", "checkout", branch)
	} else {
		return nil
	}
}

func CherryPick(PR *github.PullRequest) error {
	return utils.ExecCommand("git", "cherry-pick", *PR.Base.SHA + ".." + *PR.Head.SHA)
}

func Push(login string, project string, branch string) {
	utils.ExecCommand("git", "push", "--set-upstream", "https://github.com/" + login + "/" + project, branch, "--force")
}

func OpenPR(client *github.Client, ctx context.Context, login string, project string, head string) *github.PullRequest {
	title := "cherry-pick-bot with a bunch of commits"
	base := "master"

	open_PR, _, err := client.PullRequests.Create(
		ctx, login, project,
		&github.NewPullRequest{Title: &title, Head: &head, Base: &base})
	utils.Die(err)

	return open_PR
}

func Rebase(branch string) error {
	return utils.ExecCommand("git", "pull", "--rebase", "origin", branch)
}

func ChangeRepo(login string, project string) {
	os.MkdirAll(config.WorkDir + login, 0775)
	os.Chdir(config.WorkDir + login)
	utils.ExecCommand("git", "clone", "git://github.com/" + login + "/" + project)
	os.Chdir(config.WorkDir + login + "/" + project)
}
