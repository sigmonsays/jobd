package git

import (
	"log/slog"
	"os"
	"os/exec"
)

func DefaultGitOptions() *GitOptions {
	return &GitOptions{}
}

type GitOptions struct {
	IdentityFile string
}

func CloneRepo(opts *GitOptions, remote, gitDir string) error {
	cmdline := []string{
		"git",
		"clone",
		"-q",
		remote,
		gitDir,
	}
	slog.Info("git clone", "cmdline", cmdline)

	env := make([]string, 0)
	env = populateEnv(env, opts.IdentityFile)

	cmd := exec.Command(cmdline[0], cmdline[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	err := cmd.Run()
	if err != nil {
		return err
	}
	slog.Debug("finished git clone", "remote", remote)
	return nil
}

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

func CloneRepoWithHash(opts *GitOptions, remote string, refSpec string, gitDir string) error {

	// clone the repo
	cmdline := []string{
		"git",
		"clone",
		"-q",
		remote,
		".",
	}
	slog.Info("git clone", "cmdline", cmdline)

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
	slog.Debug("finished git clone", "remote", remote)

	// TODO: clone the specific refSpec (hash)

	return nil
}
