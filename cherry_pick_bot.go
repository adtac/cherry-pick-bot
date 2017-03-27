package main

import (
	"time"

	"github.com/google/go-github/github"

	"config"
	"utils"
	"authenticate"
	"gh"
	"replies"
	"git"
)



func main() {
	utils.LoadEnvironment()
	config.WorkDir = utils.SanitizeWorkDir(config.WorkDir)

	ctx, client := authenticate.Authenticate()

	for true {
		notifications, resp, err := client.Activity.ListNotifications(
			ctx, &github.NotificationListOptions{All: true})

		if resp.Response.StatusCode != 200 {
			utils.Die(err)
		}

		unreadNotifications := make([]*github.Notification, 0)
		for _, notification := range(notifications) {
			if notification.GetUnread() {
				unreadNotifications = append(unreadNotifications, notification)
			}
		}

		client.Activity.MarkNotificationsRead(ctx, time.Now())
		
		for _, notification := range(unreadNotifications) {
			login, project, PR_ID := gh.ExtractNotification(notification)

			git.ChangeRepo(login, project)

			if *notification.Reason == "mention" {
				// check if email is public
				mentions := gh.GetMentions(client, ctx, login, project, PR_ID)
				last_user, _, err := client.Users.Get(ctx, *mentions[len(mentions)-1].User.Login)
				utils.Die(err)
				if last_user.Email == nil {
					gh.Comment(client, ctx, login, project, PR_ID, replies.InvalidEmail)
					continue
				}

				open_PR := gh.GetOpenPullRequest(client, ctx, login, project)

				// spoof the cherry-pick committer to make it look like the person commenting
				// did it; also clear any ongoing rebases or cherry-picks
				git.SpoofUser(last_user)
				git.Clear()

				// fetch the PR (the branch actually)
				PR := gh.GetPullRequest(client, ctx, login, project, PR_ID)
				git.Fetch(PR)

				// cherry-pick the PR's commits in a new branch
				git.CheckoutBranch("cherry-pick-bot/patch")
				if git.CherryPick(PR) != nil {
					gh.Comment(client, ctx, login, project, PR_ID, replies.CannotCherryPick)
					continue
				}

				// push to github
				git.Push(login, project, "cherry-pick-bot/patch")

				comment := ""
				if open_PR == nil {
					open_PR = git.OpenPR(client, ctx, login, project, "cherry-pick-bot/patch")
					comment = "Done! Opened a new PR at " + *open_PR.HTMLURL
				} else {
					comment = "Done! Updated " + *open_PR.HTMLURL
				}

				gh.Comment(client, ctx, login, project, PR_ID, comment)
			}
		}

		time.Sleep(config.SleepTime)
	}
}
