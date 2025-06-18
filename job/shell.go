package job

import (
	"log/slog"
)

func RunShell(runSpec *RunSpec, step string, job *JobSpec, sh *ShellSpec) error {
	ex, err := NewExec(sh.Script)
	if err != nil {
		return err
	}

	// todo ex.Env
	// todo timeout
	ex.FileLog = runSpec.LogFile

	ex.ScriptDir = job.Directory
	ex.LogDir = runSpec.RunDir
	err = ex.Init()
	if err != nil {
		return err
	}

	err = ex.Run()
	if err != nil {
		return err
	}

	slog.Debug("exec returned", "jid", runSpec.JobId, "rid", runSpec.RunId, "step", step, "exitcode", ex.ExitCode)
	return nil
}
