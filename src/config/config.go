package config

import (
    "time"
)

// If the `GITHUB_ACCESS_TOKEN` environment variable is set, it would
// be used as the access token. If not, you can manually set it below.
var AccessToken = "Github personal access token"

var WorkDir = "/tmp/work_dir/"

// If the `GITHUB_EMAIL` environment variable is set, it would
// be used. Otherwise, you can manually set it below.
var Email = "email@example.com"

// If the `GITHUB_PRIVATE_KEY` environment variable is set, it would
// be used. Otherwise, you can manually set it below.
var PrivateKey = "/path/to/ssh/key"

var SleepTime = 15 * time.Second

var GithubLogin = "cherry-pick-bot"
