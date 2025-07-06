package job

import (
	"fmt"
	"io"
	"log/slog"
	"path/filepath"

	"github.com/sigmonsays/jobd/util"
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

	// disable this job
	Disabled bool

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

	// variables to capture
	VarPrefix string
	Vars      []*JobVars

	// Upstream
	Upstream *UpstreamSpec
}

type UpstreamSpec struct {
	Git *GitSpec
}

type GitSpec struct {

	// ssh key
	IdentityFile string

	Remote string

	// todo: Support more than one branch, tag, globs, etc
	Branch string
}

type JobStep struct {

	// the step id, optional
	Id string

	// disable this step
	Disabled bool

	Shell *ShellSpec
	Vars  []*JobVars
}

type ShellSpec struct {

	// the working directory
	WorkingDir string

	Env []string

	// the command to run
	Script string

	// Exit immediately if command exits non-zero (set -e)
	ImmediateExit bool

	// execution timeout
	Timeout int
}

// the run of a specific job
type RunSpec struct {
	JobId       string
	RunId       int
	RunDir      string
	WorkDir     string
	UpstreamDir string
	LogFile     io.Writer `json:"-"`

	EnvFile string

	StepResults []*StepResult

	Vars *Vars
}

func (me *RunSpec) GetLogFile() string {
	logfile := filepath.Join(me.RunDir, "job.log")
	return logfile
}

func (me *RunSpec) Logf(s string, args ...any) {
	prefix := "#### "
	fmt.Fprintf(me.LogFile, prefix+s+"\n", args...)
	msg := fmt.Sprintf(prefix+s, args...)
	slog.Debug("Logf", "msg", msg)
}

type StepResult struct {
	Id          string // step id
	Index       int
	Step        *JobStep
	Error       string
	DurationSec int

	duration *util.DurationMeasure
}

type JobVars struct {
	FromCommand []CommandVarSpec `yaml:"from_command"`
}

type CommandVarSpec map[string]string
