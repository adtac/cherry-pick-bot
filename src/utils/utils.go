package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"strconv"

	"github.com/google/go-github/github"
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

func ExtractNotification(notification *github.Notification) (login string, project string, repo string, cloneURL string, PR_ID int) {
	project = *notification.Repository.Name
	repo = *notification.Repository.FullName

	splits := strings.Split(*notification.Subject.URL, "/")
	PR_ID, err := strconv.Atoi(splits[len(splits)-1])

	Die(err)

	// git clone cloneURL
	cloneURL = "git" + (*notification.Repository.HTMLURL)[5:] + ".git"

	// username
	login = *notification.Repository.Owner.Login

	return
}
