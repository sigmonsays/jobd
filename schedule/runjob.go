package schedule

import (
	"log/slog"

	"github.com/sigmonsays/jobd/job"
)

// helper function to make a job spec runnable and execute it
func RunJobSpec(executor *Executor, j *job.JobSpec) {
	runnable := &job.RunJob{JobSpec: j}
	opts := DefaultExecuteOptions()
	opts.Stackable = j.Stackable
	err := executor.Execute(j.JobId, runnable, opts)
	if err != nil {
		slog.Warn("Execute job error", "err", err)
	}

}
