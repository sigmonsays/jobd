package polling

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sigmonsays/jobd/core"
	"github.com/sigmonsays/jobd/job"
	"github.com/sigmonsays/jobd/schedule"

	git_watch "github.com/sigmonsays/git-watch/watch/git"
)

type GitUpstreamNotify struct {
	LocalHash  string
	RemoteHash string
}

// job poller func for each job
func JobPoller(j *job.JobSpec, appCtx *core.Context) error {
	// todo: Make polling interval configurable
	pollInt := 30

	// clone repo first if needed
	gitDir := filepath.Join(j.Directory, "upstream/git/default")
	baseDir := filepath.Base(gitDir)
	os.MkdirAll(baseDir, 0766)
	_, err := os.Stat(gitDir)
	setupDir := err != nil && os.IsNotExist(err)
	if setupDir {
		err := CloneRepo(j, gitDir)
		if err != nil {
			slog.Warn("git clone error", "error", err)
		}
	}

	// start watching
	gitWatch := git_watch.NewGitWatch(gitDir, j.Upstream.Git.Branch)
	gitWatch.Interval = pollInt
	changes := make(chan *GitUpstreamNotify, 5)
	gitWatch.OnChange = func(dir, branch, lhash, rhash string) error {
		changes <- &GitUpstreamNotify{
			LocalHash:  lhash,
			RemoteHash: rhash,
		}
		return nil
	}
	gitWatch.OnCheck = func(dir, branch, lhash, rhash string) error {
		slog.Debug("Check upstream for changes", "dir", dir)
		return nil
	}
	gitWatch.Start()
	defer gitWatch.Stop()

	// todo: Wire up job schedule stopping (JobCtx) here
	for {
		select {
		case change := <-changes:

			PullRepo(j, gitDir)

			// schedule job
			// todo: Pass change event into job somehow
			slog.Info("change detected, executing job", "jid", j.JobId, "change", change)
			go schedule.RunJobSpec(appCtx.Executor, j)
		}
	}

	return nil
}
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
