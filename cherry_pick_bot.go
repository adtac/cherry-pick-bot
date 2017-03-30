package main

import (
	"time"
)

func main() {
	loadEnvironment()
	workDir = sanitizeWorkDir(workDir)

	ctx, client := authenticate()

	for true {
		unreadNotifications := getUnreadNotifications(client, ctx)

		client.Activity.MarkNotificationsRead(ctx, time.Now())
		
		for _, notification := range(unreadNotifications) {
			login, project, PR_ID := extractNotification(notification)

			changeRepo(login, project)

			if *notification.Reason == "mention" {
				// check if email is public
				last_user, err := getLastUserMentioned(client, ctx, login, project, PR_ID)

				if err != nil {
					die(err)
				}

				if last_user.Email == nil {
					comment(client, ctx, login, project, PR_ID, invalidEmail)
					continue
				}

				// spoof the cherry-pick committer to make it look like the person commenting
				// did it; also clear any ongoing rebases or cherry-picks
				spoofUser(last_user)
				clear()

				err = performCherryPick(client, ctx, login, project, PR_ID)
				if err != nil {
					continue
				}

				createCherryPR(client, ctx, login, project, PR_ID)
			}
		}

		time.Sleep(sleepTime)
	}
}
