package polling

import (
	"log/slog"
	"os"
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
	gitDir := filepath.Join(j.Directory, "upstream")
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

	// get a handle on the jobs scheduler context
	shedJob, err := appCtx.Scheduler.FindJobByName(j.JobId)
	if err != nil {
		return err
	}

	quit := shedJob.JobCtx.Done()

	// start watching
	// git is the only backend
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
Dance:
	for {
		select {
		case change := <-changes:
			go RepoChange(appCtx, j, gitDir, change)

		case <-quit:
			break Dance
		}
	}
	return nil
}

func RepoChange(appCtx *core.Context, j *job.JobSpec, gitDir string, change *GitUpstreamNotify) error {

	// pull git repo
	PullRepo(j, gitDir)

	// run job
	// todo: Pass change event into job somehow
	slog.Info("change detected, executing job", "jid", j.JobId, "change", change)
	schedule.RunJobSpec(appCtx.Executor, j)

	return nil
}
