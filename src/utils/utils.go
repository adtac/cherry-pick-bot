package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"config"
)

func LoadEnvironment() {
	m := make(map[string]*string)
	m["GITHUB_ACCESS_TOKEN"] = &config.AccessToken
	m["GITHUB_EMAIL"] = &config.Email
	m["PRIVATE_KEY"] = &config.PrivateKey
				
	for key, val := range(m) {
		var_val, present := os.LookupEnv(key)
		if present {
			*val = var_val
		}
	}

	os.Setenv("GIT_SSH_COMMAND", "ssh -i " + config.PrivateKey)
}

// sanitizes the work directory (adds a slashes at the end) and creates
// the directory
func SanitizeWorkDir(dir string) string {
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}

	os.MkdirAll(dir, 0775)

	return dir
}

func ExecCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	return cmd.Run()
}

func Die(err error) {
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}
