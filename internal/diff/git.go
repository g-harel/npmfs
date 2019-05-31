package diff

import (
	"fmt"
	"os/exec"
	"strings"
)

func execGit(dir string, arg ...string) (string, error) {
	cmd := exec.Command("git", arg...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("run command 'git %v': %v\n%v", strings.Join(arg, " "), err, string(out))
	}

	return string(out), nil
}
