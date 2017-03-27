package config

import (
    "time"
)

var ACCESS_TOKEN = "Your Github personal access token goes here."

// this HAS to end with a slash
var WORK_DIR = "/tmp/work_dir/"

var EMAIL = "email@example.com"

var PRIVATE_KEY = "Link to your SSH private key."

var SLEEP_TIME = 15 * time.seconds
