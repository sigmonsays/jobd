package git

import (
	"log/slog"
	"os"
	"os/exec"
)

func PullRepo(opts *GitOptions, gitDir string) error {
	cmdline := []string{
		"git",
		"pull",
		"-q",
		"--no-edit",
	}
	slog.Info("git pull", "cmdline", cmdline)
	env := make([]string, 0)
	env = populateEnv(env, opts.IdentityFile)

	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	cmd.Dir = gitDir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	err := cmd.Run()
	if err != nil {
		return err
	}
	slog.Debug("finished git pull", "directory", gitDir)
	return nil
}
