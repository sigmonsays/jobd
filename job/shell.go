package job

import (
	"bufio"
	"log/slog"
	"os"
	"strings"
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

	// pass through envfile
	f, err := os.Open(runSpec.EnvFile)
	if err != nil {
		return ret, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r\n")
		idx := strings.Index(line, "=")
		if idx == -1 {
			continue
		}
		k := line[:idx]
		v := line[idx+1:]
		slog.Debug("Envfile setenv", "k", k, "v", v)
		ex.SetEnv(k, v)
	}

	// pass through some global variables to shell env variables
	for k, v := range runSpec.Vars.Global.Vars {
		ex.SetEnv(k, v.String())
	}

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

	slog.Debug("exec returned",
		"jid", runSpec.JobId,
		"rid", runSpec.RunId,
		"step", step,
		"exitcode", ex.ExitCode,
		"timedout", ex.TimedOut,
	)

	err = captureVariables(job.Vars, runSpec, runSpec.Vars.GetStep(step))
	if err != nil {
		return nil, err
	}

	return ret, nil
}
