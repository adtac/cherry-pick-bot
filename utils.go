package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func loadEnvironment() {
	m := make(map[string]*string)
	m["GITHUB_ACCESS_TOKEN"] = &accessToken
	m["GITHUB_EMAIL"] = &email
	m["PRIVATE_KEY"] = &privateKey
				
	for key, val := range(m) {
		var_val, present := os.LookupEnv(key)
		if present {
			*val = var_val
		}
	}

	os.Setenv("GIT_SSH_COMMAND", "ssh -i " + privateKey)
}

// sanitizes the work directory (adds a slashes at the end) and creates
// the directory
func sanitizeWorkDir(dir string) string {
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}

	os.MkdirAll(dir, 0775)

	return dir
}

func execCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	return cmd.Run()
}

func die(err error) {
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}
