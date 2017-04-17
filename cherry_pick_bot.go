package main

import (
	"fmt"
	"flag"
	"time"

	"github.com/op/go-logging"
)

var configPath = flag.String("config", "config.toml", "Path for the config file")

var logger = logging.MustGetLogger("cherry-pick-bot")

func main() {
	err := loadConfig(*configPath)
	if err != nil {
		die(fmt.Errorf("error while reading configuration file: %v", err))
	}

	loadEnvironment()
	ctx, client := authenticate()

	logger.Notice("Ready for action!")

	for true {
		unreadNotifications, err := getUnreadNotifications(client, ctx)
		if err != nil {
			die(fmt.Errorf("error while getting unread notifications: %v", err))
		}

		logger.Infof("Got %d notifications!", len(unreadNotifications))

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

				logger.Infof("Got a call from %s on %s/%s #%d", lastUser.Login, login, project, prId)

				if lastUser.Email == nil {
					logger.Infof("%s email isn't public.. Skipping...", lastUser.Login)
					comment(client, ctx, login, project, prId, invalidEmail)
					continue
				}

				// spoof the cherry-pick committer to make it look like the person commenting
				// did it; also clear any ongoing rebases or cherry-picks
				spoofUser(lastUser)
				clear()

				logger.Infof("Performing cherry pick for %s/%s #%d ...", login, project, prId)
				err = performCherryPick(client, ctx, login, project, prId)
				if err != nil {
					logger.Error(err)
					continue
				}

				logger.Info("Creating pull request ...", login, project, prId)
				err = createCherryPR(client, ctx, login, project, prId)
				if err != nil {
					logger.Error(err)
					continue
				}
			}
		}

		logger.Info("Sleeping ...")
		time.Sleep(time.Duration(conf.SleepTime.Nanoseconds()))
	}
}
