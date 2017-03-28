package main

import (
	"time"

	"github.com/google/go-github/github"
)



func main() {
	loadEnvironment()
	workDir = sanitizeWorkDir(workDir)

	ctx, client := authenticate()

	for true {
		notifications, resp, err := client.Activity.ListNotifications(
			ctx, &github.NotificationListOptions{All: true})

		if resp.Response.StatusCode != 200 {
			die(err)
		}

		unreadNotifications := make([]*github.Notification, 0)
		for _, notification := range(notifications) {
			if notification.GetUnread() {
				unreadNotifications = append(unreadNotifications, notification)
			}
		}

		client.Activity.MarkNotificationsRead(ctx, time.Now())
		
		for _, notification := range(unreadNotifications) {
			login, project, PR_ID := extractNotification(notification)

			changeRepo(login, project)

			if *notification.Reason == "mention" {
				// check if email is public
				mentions := getMentions(client, ctx, login, project, PR_ID)
				last_user, _, err := client.Users.Get(ctx, *mentions[len(mentions)-1].User.Login)
				die(err)
				if last_user.Email == nil {
					comment(client, ctx, login, project, PR_ID, invalidEmail)
					continue
				}

				open_PR := getOpenPullRequest(client, ctx, login, project)

				// spoof the cherry-pick committer to make it look like the person commenting
				// did it; also clear any ongoing rebases or cherry-picks
				spoofUser(last_user)
				clear()

				// fetch the PR (the branch actually)
				PR := getPullRequest(client, ctx, login, project, PR_ID)
				fetch(PR)

				// cherry-pick the PR's commits in a new branch
				checkoutBranch("cherry-pick-bot/patch")
				if cherryPick(PR) != nil {
					comment(client, ctx, login, project, PR_ID, cannotCherryPick)
					continue
				}

				// push to github
				push(login, project, "cherry-pick-bot/patch")

				commentText := ""
				if open_PR == nil {
					open_PR = openPR(client, ctx, login, project, "cherry-pick-bot/patch")
					commentText = "Done! Opened a new PR at " + *open_PR.HTMLURL
				} else {
					commentText = "Done! Updated " + *open_PR.HTMLURL
				}

				comment(client, ctx, login, project, PR_ID, commentText)
			}
		}

		time.Sleep(sleepTime)
	}
}
