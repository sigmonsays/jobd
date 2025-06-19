package job

import (
	"log/slog"
)

type ShellResult struct {
	ExitCode int
	TimedOut bool
	// todo: Return more things
}

func RunShell(runSpec *RunSpec, step string, job *JobSpec, sh *ShellSpec) (*ShellResult, error) {

	slog.Debug("RunShell", "jid", job.JobId, "rid", runSpec.RunId, "dir", sh.WorkingDir)
	ex, err := NewExec(sh.Script)
	if err != nil {
		return nil, err
	}

	ret := &ShellResult{}

	ex.FileLog = runSpec.LogFile
	ex.Timeout = sh.Timeout
	ex.Dir = sh.WorkingDir
	ex.Env = sh.Env

	ex.ScriptDir = job.Directory
	ex.LogDir = runSpec.RunDir
	err = ex.Init()
	if err != nil {
		return ret, err
	}

	err = ex.Run()

	ret.ExitCode = ex.ExitCode
	ret.TimedOut = ex.TimedOut

	if err != nil {
		return ret, err
	}

	slog.Debug("exec returned", "jid", runSpec.JobId, "rid", runSpec.RunId, "step", step, "exitcode", ex.ExitCode)
	return ret, nil
}
