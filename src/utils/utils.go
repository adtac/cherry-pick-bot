package utils

import (
    "os"
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
