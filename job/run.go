package job

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type RunJob struct {
	JobSpec *JobSpec
}

func (me *RunJob) Run() (*RunSpec, error) {
	return Run(me.JobSpec)
}

func Run(job *JobSpec) (*RunSpec, error) {
	rundir_base := filepath.Join(job.Directory, "run")
	runid, err := getMaxRun(rundir_base)
	if err != nil {
		return nil, err
	}
	runid = runid + 1
	runid_str := fmt.Sprintf("%d", runid)
	rundir := filepath.Join(rundir_base, runid_str, "")
	workdir := filepath.Join(rundir, "workspace")
	os.MkdirAll(rundir, 0700)
	os.MkdirAll(workdir, 0700)

	// make vars api

	vars := NewVars(job.VarPrefix)
	vars.Global.SetVar("JOBID", job.JobId)

	// build run spe
	runSpec := &RunSpec{
		JobId:   job.JobId,
		RunId:   runid,
		RunDir:  rundir,
		WorkDir: workdir,
		Vars:    vars,
	}

	envfile := filepath.Join(runSpec.RunDir, "shell.env")

	vars.Global.SetVar("RUNID", runid_str)
	vars.Global.SetVar("RUNDIR", rundir)
	vars.Global.SetVar("WORKSPACE", workdir)
	vars.Global.SetVar("ENV", envfile)

	slog.Debug("run job", "rundir", rundir, "workdir", workdir)

	if len(job.Steps) == 0 {
		return nil, fmt.Errorf("Job has no steps: jid:%s", job.JobId)
	}

	// open output file for job
	open_flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	stdoutpath := filepath.Join(runSpec.RunDir, "job.log")
	outfh, err := os.OpenFile(stdoutpath, open_flags, 0644)
	if err != nil {
		return nil, err
	}
	runSpec.LogFile = outfh

	defer outfh.Close()

	// log line for job
	runSpec.Logf("run job jid:%s runid:%d",
		job.JobId, runSpec.RunId)

	for step_idx, step := range job.Steps {

		if step.Disabled {
			runSpec.Logf("run step jid:%s step_num:%d: name:%s: step is disabled, skipping",
				job.JobId, step_idx, step.Id)
		}

		stepRes := &StepResult{
			Id:    step.Id,
			Index: step_idx,
			Step:  step,
		}
		runSpec.StepResults = append(runSpec.StepResults, stepRes)

		runSpec.Logf("run step jid:%s step_num:%d: name:%s",
			job.JobId, step_idx, step.Id)

		job_type := "unknown"
		if step.Shell != nil {
			job_type = "shell"
		}

		if job_type == "unknown" {
			slog.Warn("Unknown job type", "job_type", job_type)
			stepRes.Error = "unknown job type"
			continue
		}

		if job_type == "shell" {

			if step.Shell.WorkingDir == "" {
				step.Shell.WorkingDir = runSpec.WorkDir
			}

			// run shell
			shres, err := RunShell(runSpec, step.Id, job, step.Shell)
			if err == nil {
				runSpec.Logf("run job jid:%s runid:%d finished, exited %d, timedout %v",
					job.JobId, runSpec.RunId, shres.ExitCode, shres.TimedOut)
			} else {
				slog.Warn("RunShell", "error", err)
				stepRes.Error = err.Error()

				var ee *exec.ExitError
				if errors.As(err, &ee) {
					runSpec.Logf("run job jid:%s runid:%d: exit error %s",
						job.JobId, runSpec.RunId, ee)

				} else {
					runSpec.Logf("run job jid:%s runid:%d: generic error %s",
						job.JobId, runSpec.RunId, err)
				}
			}
		}

	}

	err = captureVariables(job.Vars, runSpec, runSpec.Vars.Global)
	if err != nil {
		return nil, err
	}

	// purge runs
	if job.Keep > 0 {
		err = purgeRuns(rundir_base, job.Keep)
		if err != nil {
			return nil, err
		}
	}

	return runSpec, nil
}

func captureVariables(vars_spec []*JobVars, runSpec *RunSpec, vars *VarSet) error {

	// capture variables
	for _, jv := range vars_spec {

		// process from command variables
		for _, command_vars := range jv.FromCommand {

			for var_name, command := range command_vars {

				cmdline := []string{
					"sh", "-x", "-c", command,
				}
				c := exec.Command(cmdline[0], cmdline[1:]...)
				buf := bytes.NewBuffer(nil)
				c.Stderr = os.Stderr
				c.Stdout = buf
				c.Dir = runSpec.WorkDir
				err := c.Run()
				if err != nil {
					slog.Warn("From_command error", "variable", var_name, "command", command, "error", err)
					continue
				}
				val := strings.Trim(buf.String(), " \n\t")

				slog.Debug("command variable result", "variable", var_name, "value", val)
				vars.SetVar(var_name, val)
			}
		}

	}
	return nil
}

// return the highest integer directory in the jobs rundir
func getMaxRun(rundir string) (int, error) {
	ret := 0

	files, err := ioutil.ReadDir(rundir)
	if err != nil {
		return 0, err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		_val, err := strconv.ParseInt(file.Name(), 10, 32)
		if err != nil {
			continue
		}
		val := int(_val)
		if val > ret {
			ret = val
		}
	}

	return ret, nil
}
func purgeRuns(rundir string, keep int) error {
	files, err := ioutil.ReadDir(rundir)
	if err != nil {
		return err
	}
	names := make([]string, 0)

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		_, err := strconv.ParseInt(file.Name(), 10, 32)
		if err != nil {
			continue
		}
		names = append(names, file.Name())
	}

	sort.Slice(names, func(i, j int) bool {
		a, _ := strconv.ParseInt(names[i], 10, 32)
		b, _ := strconv.ParseInt(names[j], 10, 32)
		return a > b
	})

	// slog.Debug("considering purge on names", "names", names)

	for idx, name := range names {
		dir := filepath.Join(rundir, name)
		if idx+1 > keep {
			slog.Debug("Purge run directory", "dir", dir)
			os.RemoveAll(dir)
		}
	}
	return nil
}
