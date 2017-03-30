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
			login, project, prId := extractNotification(notification)

			changeRepo(login, project)

			if *notification.Reason == "mention" {
				// check if email is public
				lastUser, err := getLastUserMentioned(client, ctx, login, project, prId)

				if err != nil {
					die(err)
				}

				if lastUser.Email == nil {
					comment(client, ctx, login, project, prId, invalidEmail)
					continue
				}

				// spoof the cherry-pick committer to make it look like the person commenting
				// did it; also clear any ongoing rebases or cherry-picks
				spoofUser(lastUser)
				clear()

				err = performCherryPick(client, ctx, login, project, prId)
				if err != nil {
					continue
				}

				createCherryPR(client, ctx, login, project, prId)
			}
		}

		time.Sleep(sleepTime)
	}
}
