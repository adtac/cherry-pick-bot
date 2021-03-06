package main

import (
    "errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"github.com/google/go-github/github"
)

func getMentions(client *github.Client, ctx context.Context, login string, project string, prId int) ([]*github.IssueComment, error) {
	comments, _, err := client.Issues.ListComments(ctx, login, project, prId, &github.IssueListCommentsOptions{Direction: "asc"})
	if err != nil {
		return make([]*github.IssueComment, 0), err
	}

	res := make([]*github.IssueComment, 0)
	for _, comment := range(comments) {
		if strings.Contains(*comment.Body, "@" + conf.GithubLogin) {
			res = append(res, comment)
		}
	}

	return res, nil
}

func getLastUserMentioned(client *github.Client, ctx context.Context, login string, project string, prId int) (*github.User, error) {
	mentions, err := getMentions(client, ctx, login, project, prId)
	if err != nil {
		return nil, err
	}

    if len(mentions) == 0 {
        return nil, errors.New("zero mentions")
    }

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

func extractNotification(notification *github.Notification) (login string, project string, prId int, err error) {
	project = *notification.Repository.Name

	splits := strings.Split(*notification.Subject.URL, "/")
	prId, err = strconv.Atoi(splits[len(splits)-1])

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

func getUnreadNotifications(client *github.Client, ctx context.Context) ([]*github.Notification, error) {
	notifications, resp, err := client.Activity.ListNotifications(
			ctx, &github.NotificationListOptions{All: true})

	if err != nil {
		return nil, err
	} else if s := resp.Response.StatusCode; s != 200 {
		return nil, fmt.Errorf("response status code is %d", s)
	}

	unreadNotifications := make([]*github.Notification, 0)
	for _, notification := range(notifications) {
		if notification.GetUnread() {
			unreadNotifications = append(unreadNotifications, notification)
		}
	}
	return unreadNotifications, nil
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

func createCherryPR(client *github.Client, ctx context.Context, login string, project string, prId int) error {
	pr := getOpenPullRequest(client, ctx, login, project)
	commentText := ""
	if pr == nil {
		pr, err := openPR(client, ctx, login, project, "cherry-pick-bot/patch")
		if err != nil {
			return err
		}
		commentText = "Done! Opened a new PR at " + *pr.HTMLURL
	} else {
		commentText = "Done! Updated " + *pr.HTMLURL
	}

	comment(client, ctx, login, project, prId, commentText)

	return nil
}
