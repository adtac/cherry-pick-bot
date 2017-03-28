package main

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"github.com/google/go-github/github"
)

func GetMentions(client *github.Client, ctx context.Context, login string, project string, PR_ID int) []*github.IssueComment {
	comments, _, err := client.Issues.ListComments(ctx, login, project, PR_ID, &github.IssueListCommentsOptions{Direction: "asc"})
	Die(err)

	res := make([]*github.IssueComment, 0)
	for _, comment := range(comments) {
		if strings.Contains(*comment.Body, "@" + GithubLogin) {
			res = append(res, comment)
		}
	}

	return res
}

func Comment(client *github.Client, ctx context.Context, login string, project string, PR_ID int, c string) {
	comment := &github.IssueComment{Body: &c}
	client.Issues.CreateComment(ctx, login, project, PR_ID, comment)
}

func GetPullRequest(client *github.Client, ctx context.Context, login string, project string, PR_ID int) *github.PullRequest {
	u := "/repos/" + login + "/" + project + "/pulls/" + fmt.Sprintf("%d", PR_ID)
	req, _ := client.NewRequest("GET", u, nil)
	pull := new(github.PullRequest)
	client.Do(ctx, req, pull)

	return pull
}

func ExtractNotification(notification *github.Notification) (login string, project string, PR_ID int) {
	project = *notification.Repository.Name

	splits := strings.Split(*notification.Subject.URL, "/")
	PR_ID, err := strconv.Atoi(splits[len(splits)-1])
	Die(err)

	login = *notification.Repository.Owner.Login

	return
}

func GetOpenPullRequest(client *github.Client, ctx context.Context, login string, project string) (PR *github.PullRequest) {
	prs, _, _ := client.PullRequests.List(ctx, login, project, &github.PullRequestListOptions{})
	for _, pr := range(prs) {
		if *pr.User.Login == "cherry-pick-bot" {
			PR = pr
		}
	}

	return
}
