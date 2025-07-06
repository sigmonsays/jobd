package schedule

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/sigmonsays/jobd/job"
)

func DefaultExecutorOptions() *ExecutorOptions {
	return &ExecutorOptions{}
}

type ExecutorOptions struct {
	JobDir string
}

func NewExecutor(opts *ExecutorOptions) *Executor {
	return &Executor{
		opts:    opts,
		results: make(map[string]*Result, 0),
		running: make(map[string]bool, 0),
	}
}

type Executor struct {
	opts *ExecutorOptions

	claimLock sync.Mutex
	mx        sync.Mutex

	running map[string]bool

	results map[string]*Result
}

type Runnable interface {
	Run() (*job.RunSpec, error)
}

func (me *Executor) unclaimJid(jid string) error {
	me.claimLock.Lock()
	defer me.claimLock.Unlock()
	delete(me.running, jid)
	return nil
}

func (me *Executor) claimJid(jid string) error {
	me.claimLock.Lock()
	defer me.claimLock.Unlock()

	_, found := me.running[jid]
	if found {
		return fmt.Errorf("already running")
	}

	me.running[jid] = true
	return nil
}

func DefaultExecuteOptions() *ExecuteOptions {
	return &ExecuteOptions{}
}

type ExecuteOptions struct {
	Stackable bool
}

func (me *Executor) Execute(jid string, f Runnable, opts *ExecuteOptions) error {
	started := time.Now()

	if opts.Stackable == false {
		// claim this slot so it only runs once
		err := me.claimJid(jid)
		if err != nil {
			slog.Debug("unable to claim job, already running", "jid", jid)
			return err
		}
		defer me.unclaimJid(jid)
	}

	// begin executing job
	slog.Debug("Execute job", "jid", jid)
	runSpec, err := f.Run()
	res := &Result{
		Error:         err,
		ExitCode:      0,
		StopTimestamp: time.Now(),
		RunSpec:       runSpec,
	}
	res.StartTimestamp = started
	res.Duration = res.StopTimestamp.Sub(res.StartTimestamp)
	res.DurationSec = int(res.Duration.Seconds())

	// extract exit code
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		res.ExitCode = ee.ExitCode()
	}

	slog.Debug("Execute job return", "jid", jid, "result", res)

	// save results on disk
	resultFile := filepath.Join(runSpec.RunDir, "result.json")
	rbuf, _ := json.Marshal(res)
	os.WriteFile(resultFile, rbuf, 0644)
	slog.Debug("saved results file", "jid", jid, "result_file", resultFile)

	// save results ni memory
	me.mx.Lock()
	me.results[jid] = res
	me.mx.Unlock()

	return err
}

func (me *Executor) GetResult(jid string) (*Result, error) {

	res, err := me.GetResultMemory(jid)
	if err != nil {
		var ee *ExecutorError
		if errors.As(err, &ee) {
			if ee.Reason == "NOT_FOUND" {
				slog.Debug("GetResult ExecutorError", "jid", jid, "ee", ee)
			} else {
				return nil, err
			}
		}
	}

	if res != nil {
		return res, nil
	}

	// try file system
	slog.Debug("GetResult, trying filesystem", "jid", jid)
	jobDir := filepath.Join(me.opts.JobDir, jid)
	rundir := filepath.Join(jobDir, "run")
	maxrun, err := job.GetMaxRun(rundir)
	runid := fmt.Sprintf("%d", maxrun)
	resultFile := filepath.Join(rundir, runid, "result.json")
	buf, err := os.ReadFile(resultFile)
	if err == nil {
		json.Unmarshal(buf, &res)
		slog.Debug("Loaded job result from disk", "result_file", resultFile)
	} else {
		slog.Debug("Error loading job result from disk", "result_file", resultFile, "err", err)

	}

	return res, err
}

type ExecutorError struct {
	Reason string
}

func (me *ExecutorError) Error() string {
	return "ExecutorError:" + me.Reason
}

func (me *Executor) GetResultMemory(jid string) (*Result, error) {
	me.mx.Lock()
	defer me.mx.Unlock()
	res, found := me.results[jid]
	if !found {
		e := &ExecutorError{"NOT_FOUND"}
		return nil, fmt.Errorf("executor jid not found: %w", e)
	}
	return res, nil
}

type Result struct {
	Error    error
	ExitCode int

	RunSpec *job.RunSpec

	StopTimestamp  time.Time
	StartTimestamp time.Time
	Duration       time.Duration
	DurationSec    int
}
