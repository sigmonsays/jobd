package polling

import (
	"log/slog"
	"os"
	"os/exec"

	"github.com/sigmonsays/jobd/job"
)

func CloneRepo(j *job.JobSpec, gitDir string) error {

	cmdline := []string{
		"git",
		"clone",
		"-q",
		j.Upstream.Git.Remote,
		gitDir,
	}
	slog.Info("git clone", "cmdline", cmdline)

	env := make([]string, 0)
	env = populateEnv(env, j)

	cmd := exec.Command(cmdline[0], cmdline[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func PullRepo(j *job.JobSpec, gitDir string) error {

	cmdline := []string{
		"git",
		"pull",
		"-q",
		"--no-edit",
	}
	slog.Info("git pull", "cmdline", cmdline)
	env := make([]string, 0)
	env = populateEnv(env, j)

	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	cmd.Dir = gitDir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
