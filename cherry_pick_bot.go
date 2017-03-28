package main

import (
	"time"

	"github.com/google/go-github/github"
)



func main() {
	LoadEnvironment()
	WorkDir = SanitizeWorkDir(WorkDir)

	ctx, client := Authenticate()

	for true {
		notifications, resp, err := client.Activity.ListNotifications(
			ctx, &github.NotificationListOptions{All: true})

		if resp.Response.StatusCode != 200 {
			Die(err)
		}

		unreadNotifications := make([]*github.Notification, 0)
		for _, notification := range(notifications) {
			if notification.GetUnread() {
				unreadNotifications = append(unreadNotifications, notification)
			}
		}

		client.Activity.MarkNotificationsRead(ctx, time.Now())
		
		for _, notification := range(unreadNotifications) {
			login, project, PR_ID := ExtractNotification(notification)

			ChangeRepo(login, project)

			if *notification.Reason == "mention" {
				// check if email is public
				mentions := GetMentions(client, ctx, login, project, PR_ID)
				last_user, _, err := client.Users.Get(ctx, *mentions[len(mentions)-1].User.Login)
				Die(err)
				if last_user.Email == nil {
					Comment(client, ctx, login, project, PR_ID, InvalidEmail)
					continue
				}

				open_PR := GetOpenPullRequest(client, ctx, login, project)

				// spoof the cherry-pick committer to make it look like the person commenting
				// did it; also clear any ongoing rebases or cherry-picks
				SpoofUser(last_user)
				Clear()

				// fetch the PR (the branch actually)
				PR := GetPullRequest(client, ctx, login, project, PR_ID)
				Fetch(PR)

				// cherry-pick the PR's commits in a new branch
				CheckoutBranch("cherry-pick-bot/patch")
				if CherryPick(PR) != nil {
					Comment(client, ctx, login, project, PR_ID, CannotCherryPick)
					continue
				}

				// push to github
				Push(login, project, "cherry-pick-bot/patch")

				comment := ""
				if open_PR == nil {
					open_PR = OpenPR(client, ctx, login, project, "cherry-pick-bot/patch")
					comment = "Done! Opened a new PR at " + *open_PR.HTMLURL
				} else {
					comment = "Done! Updated " + *open_PR.HTMLURL
				}

				Comment(client, ctx, login, project, PR_ID, comment)
			}
		}

		time.Sleep(SleepTime)
	}
}
