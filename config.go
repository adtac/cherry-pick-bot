package main

import (
    "time"
)

// If the `GITHUB_ACCESS_TOKEN` environment variable is set, it would
// be used as the access token. If not, you can manually set it below.
var accessToken = "Github personal access token"

var workDir = "/tmp/work_dir/"

// If the `GITHUB_EMAIL` environment variable is set, it would
// be used. Otherwise, you can manually set it below.
var email = "email@example.com"

// If the `GITHUB_PRIVATE_KEY` environment variable is set, it would
// be used. Otherwise, you can manually set it below.
var privateKey = "/path/to/ssh/key"

var sleepTime = 15 * time.Second

var githubLogin = "cherry-pick-bot"
