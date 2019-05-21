package git

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func execGit(dir string, arg ...string) (io.Reader, error) {
	cmd := exec.Command("git", arg...)
	cmd.Dir = dir

	// TODO use reader/writer to reduce potential memory usage.
	// https://golang.org/pkg/os/exec/#example_Cmd_StdoutPipe
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("run command 'git %v': %v\n%v", strings.Join(arg, " "), err, string(out))
	}

	return bytes.NewReader(out), nil
}
