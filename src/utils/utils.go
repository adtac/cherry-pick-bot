package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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
	panic(fmt.Sprintf("%s", err))
}
