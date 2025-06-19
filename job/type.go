package job

import (
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
)

/*
  JID: String - job name
  RID: int - run id

  datadir /srv/jobd
  jobdir /srv/jobd/jobs/JID/
  - stdout /srv/jobd/jobs/JID/run/RID/stdout.log
  - stderr /srv/jobd/jobs/JID/run/RID/stderr.log

  run
  - /srv/jobd/jobs/JID/run/RID/job.json
*/

type JobSpec struct {

	// job id (JID)
	JobId string

	// immediately run the job
	Immediate bool

	// job params
	// stackable is default false, true means multiple can run concurrently
	Stackable bool

	// how many runs to keep
	Keep int

	// top level directory on disk of job
	Directory string

	Steps []*JobStep

	// run schedule
	Schedule string
}

type JobStep struct {

	// the step id, optional
	Id string

	Shell *ShellSpec
}

type ShellSpec struct {

	// the working directory
	WorkingDir string

	Env []string

	// the command to run
	Script string

	// execution timeout
	Timeout int
}

// the run of a specific job
type RunSpec struct {
	JobId   string
	RunId   int
	RunDir  string
	LogFile io.Writer `json:"-"`

	StepResults []*StepResult
}

func (me *RunSpec) GetLogFile() string {
	logfile := filepath.Join(me.RunDir, "job.log")
	return logfile
}

func (me *RunSpec) Logf(s string, args ...any) {
	fmt.Fprintf(me.LogFile, s, args...)
	msg := fmt.Sprintf("#### "+s+"\n", args...)
	slog.Debug("Logf", "msg", msg)
}

type StepResult struct {
	Id    string // step id
	Index int
	Step  *JobStep
	Error string
}
