package git

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

// get a commit hash for branch in git repo at directory
func LocalHash(opts *GitOptions, dir, branch string) string {
	cmdline := []string{
		"-C", dir,
		"rev-list",
		"--max-count=1",
		branch,
	}
	cmd := exec.Command("git", cmdline...)
	e := make([]string, 0)
	e = populateEnv(e, opts.IdentityFile)
	cmd.Env = e
	out, err := cmd.Output()
	if err != nil {
		slog.Info("local hash error", "cmdline", cmdline, "error", err)
		return ""
	}
	return strings.Trim(string(out), "\n")
}

// get remote hash
func RemoteHash(opts *GitOptions, dir, branch string) string {
	cmdline := []string{
		"-C", dir,
		"ls-remote",
		"origin",
		"-h",
		fmt.Sprintf("refs/heads/%s", branch),
	}
	cmd := exec.Command("git", cmdline...)
	e := make([]string, 0)
	e = populateEnv(e, opts.IdentityFile)
	cmd.Env = e

	out, err := cmd.Output()
	if err != nil {
		slog.Info("remote hash error", "cmdline", cmdline, "error", err)
		return ""
	}
	tmp := strings.Fields(string(out))
	if len(tmp) < 1 {
		return ""
	}
	return tmp[0]
}
