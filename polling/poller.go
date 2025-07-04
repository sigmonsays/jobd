package polling

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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
	pollInt := time.Duration(30) * time.Second

	tick := time.NewTicker(pollInt)
	defer tick.Stop()

	// clone repo first if needed
	gitDir := filepath.Join(j.Directory, "upstream/git")
	os.MkdirAll(gitDir, 0766)
	_, err := os.Stat(gitDir)
	setupDir := err != nil && os.IsNotExist(err)
	if setupDir {

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
			slog.Warn("git clone error", "error", err)
		}
	}

	// start watching
	gitWatch := git_watch.NewGitWatch(j.Directory, j.Upstream.Git.Branch)
	changes := make(chan *GitUpstreamNotify, 5)
	gitWatch.OnChange = func(dir, branch, lhash, rhash string) error {
		changes <- &GitUpstreamNotify{
			LocalHash:  lhash,
			RemoteHash: rhash,
		}
		return nil
	}

	for {
		select {
		case change := <-changes:

			// schedule job
			// todo: Pass change event into job somehow
			slog.Info("change detected, executing job", "jid", j.JobId, "change", change)
			go schedule.RunJobSpec(appCtx.Executor, j)
		}
	}

	return nil
}
