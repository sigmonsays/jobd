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
		j.Upstream.Git.Remote,
		gitDir,
	}
	slog.Info("git clone", "cmdline", cmdline)

	cmd := exec.Command(cmdline[0], cmdline[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
	}
	slog.Info("git pull", "cmdline", cmdline)

	cmd := exec.Command(cmdline[0], cmdline[1:]...)
	cmd.Dir = gitDir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
