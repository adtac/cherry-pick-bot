package main

import (
	"strings"
	"os"
	"fmt"
	"time"
	"strconv"

	"github.com/google/go-github/github"

	"config"
	"utils"
	"authenticate"
)



func main() {
	work_dir := utils.SanitizeWorkDir(config.WorkDir)

	os.Setenv("GIT_SSH_COMMAND", "ssh -i " + config.PrivateKey)

	ctx, client := authenticate.Authenticate(config.AccessToken)

	for true {
		notifications, resp, err := client.Activity.ListNotifications(
			ctx, &github.NotificationListOptions{All: true})

		if resp.Response.StatusCode != 200 {
			utils.Die(err)
		}
		
		for _, notification := range(notifications) {
			if notification.GetUnread() {
				if *notification.Reason == "mention" {
					project := *notification.Repository.Name
					repo := *notification.Repository.FullName
					fmt.Println("cherry-pick " + repo)

					_tmp := strings.Split(*notification.Subject.URL, "/")
					PR_ID, _ := strconv.Atoi(_tmp[len(_tmp)-1])

					cloneURL := "git" + (*notification.Repository.HTMLURL)[5:] + ".git"
					login := *notification.Repository.Owner.Login

					var PR *github.PullRequest = nil
					prs, _, _ := client.PullRequests.List(ctx, login, project, &github.PullRequestListOptions{})
					for _, pr := range(prs) {
						if *pr.User.Login == "cherry-pick-bot" {
							PR = pr
						}
					}

					os.MkdirAll(work_dir + login, 0775)
					os.Chdir(work_dir + login)
					utils.ExecCommand("git", "clone", cloneURL)
					os.Chdir(work_dir + repo)

					utils.ExecCommand("git", "config", "user.email", config.Email)
					utils.ExecCommand("git", "config", "user.name", "Cherry Pick Bot")
					utils.ExecCommand("git", "remote", "set-url", "origin", "git@github.com:" + repo + ".git")
					utils.ExecCommand("git", "cherry-pick", "--abort")
					utils.ExecCommand("git", "rebase", "--abort")

					u := (*notification.Subject.URL)[22:]
					req, _ := client.NewRequest("GET", u, nil)
					pull := new(github.PullRequest)
					client.Do(ctx, req, pull)

					creator := *pull.User.Login
					utils.ExecCommand("git", "remote", "add", creator, *pull.Head.Repo.GitURL)
					utils.ExecCommand("git", "fetch", creator)

					if utils.ExecCommand("git", "checkout", "-b", "cherry-pick-bot/patch") != nil {
						fmt.Println("branch probably exists")
						if utils.ExecCommand("git", "checkout", "cherry-pick-bot/patch") != nil {
							fmt.Println("nope, can't create/switch to branch")
							continue
						}
					}

					if utils.ExecCommand("git", "cherry-pick", *pull.Base.SHA + ".." + *pull.Head.SHA) != nil {
						c := "Uh-oh. I can't cherry-pick these commits. Any of the following could be the reason:\n\n- There are conflicts due to other commits being cherry-picked before.\n- Something has been merged into master and that's causing a conflict (in this case, ask the author of this commit to rebase to master and resolve all conflicts; nothing I can do here).\n- These commits have already been added for cherry-picking! If the commits have changed since, please close that PR and cherry-pick everything again."
						comment := &github.IssueComment{Body: &c}
						client.Issues.CreateComment(ctx, login, project, PR_ID, comment)
						fmt.Println("can't cherry-pick")
						continue
					}

					utils.ExecCommand("git", "push", "--set-upstream", "origin", "cherry-pick-bot/patch", "--force")

					if PR == nil {
						title := "cherry-pick-bot with a bunch of commits"
						head := "cherry-pick-bot/patch"
						base := "master"

						PR, _, _ = client.PullRequests.Create(
							ctx, login, project,
							&github.NewPullRequest{Title: &title, Head: &head, Base: &base})

						c := "Done! Opened a new PR at " + *PR.HTMLURL
						comment := &github.IssueComment{Body: &c}
						client.Issues.CreateComment(ctx, login, project, PR_ID, comment)
					} else {
						c := "Done! Updated " + *PR.HTMLURL
						comment := &github.IssueComment{Body: &c}
						client.Issues.CreateComment(ctx, login, project, PR_ID, comment)
					}
				} else if *notification.Reason == "author" {
					project := *notification.Repository.Name
					repo := *notification.Repository.FullName
					fmt.Println("rebase " + repo)

					_tmp := strings.Split(*notification.Subject.URL, "/")
					PR_ID, _ := strconv.Atoi(_tmp[len(_tmp)-1])

					cloneURL := "git" + (*notification.Repository.HTMLURL)[5:] + ".git"
					login := *notification.Repository.Owner.Login

					os.MkdirAll(work_dir + login, 0775)
					os.Chdir(work_dir + login)
					utils.ExecCommand("git", "clone", cloneURL)
					os.Chdir(work_dir + repo)

					utils.ExecCommand("git", "config", "user.email", config.Email)
					utils.ExecCommand("git", "config", "user.name", "Cherry Pick Bot")
					utils.ExecCommand("git", "remote", "set-url", "origin", "git@github.com:" + repo + ".git")
					utils.ExecCommand("git", "cherry-pick", "--abort")
					utils.ExecCommand("git", "rebase", "--abort")

					utils.ExecCommand("git", "checkout", "cherry-pick-bot/patch")
					if utils.ExecCommand("git", "pull", "--rebase", "origin", "master") != nil {
						c := "Uh-oh. I couldn't rebase. This may happen because master has changed a lot and there are conflicts now. I can't really resolve conflicts, so you're going to have to do this one manually. Sorry!"
						comment := &github.IssueComment{Body: &c}
						client.Issues.CreateComment(ctx, login, project, PR_ID, comment)
						fmt.Println("can't rebase")
						continue
					}

					utils.ExecCommand("git", "push", "--set-upstream", "origin", "cherry-pick-bot/patch", "--force")

					c := "Done! Rebased this PR."
					comment := &github.IssueComment{Body: &c}
					client.Issues.CreateComment(ctx, login, project, PR_ID, comment)
				}
			}
		}

		client.Activity.MarkNotificationsRead(ctx, time.Now())

		fmt.Println("sleep")
		time.Sleep(config.SleepTime)
	}
}
