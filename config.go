package main

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return
}

type config struct {
	AccessToken string
	WorkDir     string
	Email       string
	PrivateKey  string
	SleepTime   duration
	GithubLogin string
}

var conf config

func loadConfig(filename string) (err error) {
	logger.Infof("Loading configuration file (%s) ...", filename)
	_, err = toml.DecodeFile(filename, &conf)
	conf.WorkDir = sanitizeWorkDir(conf.WorkDir)
	return
}

func loadEnvironment() {
	logger.Info("Looking up environment variables ...")
	m := make(map[string]*string)
	m["GITHUB_ACCESS_TOKEN"] = &conf.AccessToken
	m["GITHUB_EMAIL"] = &conf.Email
	m["PRIVATE_KEY"] = &conf.PrivateKey

	for key, val := range(m) {
		varVal, present := os.LookupEnv(key)
		if present {
			*val = varVal
		}
	}

	os.Setenv("GIT_SSH_COMMAND", "ssh -i " + conf.PrivateKey)
}
