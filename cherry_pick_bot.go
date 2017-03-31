package main

import (
	"fmt"
	"flag"
	"time"
)

var configPath = flag.String("config", "config.toml", "Path for the config file")

func main() {
	err := loadConfig(*configPath)
	if err != nil {
		die(fmt.Errorf("error while reading configuration file: %v", err))
	}

	loadEnvironment()
	ctx, client := authenticate()

	for true {
		unreadNotifications, err := getUnreadNotifications(client, ctx)
		if err != nil {
			die(fmt.Errorf("error while getting unread notifications: %v", err))
		}

		client.Activity.MarkNotificationsRead(ctx, time.Now())
		
		for _, notification := range(unreadNotifications) {
			login, project, prId, err := extractNotification(notification)
			if err != nil {
				die(fmt.Errorf("error while extracting notification data: %v", err))
			}

			changeRepo(login, project)

			if *notification.Reason == "mention" {
				// check if email is public
				lastUser, err := getLastUserMentioned(client, ctx, login, project, prId)

				if err != nil {
					die(fmt.Errorf("error while getting mentioner: %v", err))
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

				err = createCherryPR(client, ctx, login, project, prId)
				if err != nil {
					continue
				}
			}
		}

		time.Sleep(time.Duration(conf.SleepTime.Nanoseconds()))
	}
}
