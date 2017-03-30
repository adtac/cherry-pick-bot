package main

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"github.com/google/go-github/github"
)

func getMentions(client *github.Client, ctx context.Context, login string, project string, prId int) []*github.IssueComment {
	comments, _, err := client.Issues.ListComments(ctx, login, project, prId, &github.IssueListCommentsOptions{Direction: "asc"})
	die(err)

	res := make([]*github.IssueComment, 0)
	for _, comment := range(comments) {
		if strings.Contains(*comment.Body, "@" + githubLogin) {
			res = append(res, comment)
		}
	}

	return res
}

func getLastUserMentioned(client *github.Client, ctx context.Context, login string, project string, prId int) (*github.User, error) {
	mentions := getMentions(client, ctx, login, project, prId)
	lastUser, _, err := client.Users.Get(ctx, *mentions[len(mentions)-1].User.Login)
	if err != nil {
		return nil, err
	}
	return lastUser, nil
}

func comment(client *github.Client, ctx context.Context, login string, project string, prId int, c string) {
	comment := &github.IssueComment{Body: &c}
	client.Issues.CreateComment(ctx, login, project, prId, comment)
}

func getPullRequest(client *github.Client, ctx context.Context, login string, project string, prId int) *github.PullRequest {
	u := "/repos/" + login + "/" + project + "/pulls/" + fmt.Sprintf("%d", prId)
	req, _ := client.NewRequest("GET", u, nil)
	pull := new(github.PullRequest)
	client.Do(ctx, req, pull)

	return pull
}

func extractNotification(notification *github.Notification) (login string, project string, prId int) {
	project = *notification.Repository.Name

	splits := strings.Split(*notification.Subject.URL, "/")
	prId, err := strconv.Atoi(splits[len(splits)-1])
	die(err)

	login = *notification.Repository.Owner.Login

	return
}

func getOpenPullRequest(client *github.Client, ctx context.Context, login string, project string) (result *github.PullRequest) {
	prs, _, _ := client.PullRequests.List(ctx, login, project, &github.PullRequestListOptions{})
	for _, pr := range(prs) {
		if *pr.User.Login == "cherry-pick-bot" {
			result = pr
		}
	}

	return
}

func getUnreadNotifications(client *github.Client, ctx context.Context) []*github.Notification {
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
	return unreadNotifications
}


func performCherryPick(client *github.Client, ctx context.Context, login string, project string, prId int) error {
	// fetch the PR (the branch actually)
	pr := getPullRequest(client, ctx, login, project, prId)
	fetch(pr)

	// cherry-pick the PR's commits in a new branch
	checkoutBranch("cherry-pick-bot/patch")
	if err := cherryPick(pr); err != nil {
		comment(client, ctx, login, project, prId, cannotCherryPick)
		return fmt.Errorf("cannot cherry-pick commits: %v", err)
	}

	// push to github
	push(login, project, "cherry-pick-bot/patch")

	return nil
}

func createCherryPR(client *github.Client, ctx context.Context, login string, project string, prId int) {
	pr := getOpenPullRequest(client, ctx, login, project)
	commentText := ""
	if pr == nil {
		pr = openPR(client, ctx, login, project, "cherry-pick-bot/patch")
		commentText = "Done! Opened a new PR at " + *pr.HTMLURL
	} else {
		commentText = "Done! Updated " + *pr.HTMLURL
	}

	comment(client, ctx, login, project, prId, commentText)
}
