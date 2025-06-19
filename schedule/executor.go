package schedule

import (
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"sync"
	"time"

	"github.com/sigmonsays/jobd/job"
)

func NewExecutor() *Executor {
	return &Executor{
		results: make(map[string]*Result, 0),
		running: make(map[string]bool, 0),
	}
}

type Executor struct {
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
	me.mx.Lock()
	me.results[jid] = res
	me.mx.Unlock()

	return err
}
func (me *Executor) GetResult(jid string) (*Result, error) {
	me.mx.Lock()
	defer me.mx.Unlock()
	res, found := me.results[jid]
	if !found {
		return nil, fmt.Errorf("executor jid not found:%s", jid)
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
