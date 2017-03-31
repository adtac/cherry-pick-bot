package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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
