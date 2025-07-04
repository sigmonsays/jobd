package job

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func NewExec(script string) (*Exec, error) {
	e := &Exec{
		Script: script,
	}
	return e, nil
}

type Exec struct {

	// where we log
	FileLog io.Writer

	ImmediateExit bool // exit immediately on command exit non-zero (set -e)

	// directory to enter
	Dir string

	// where we store job.sh
	ScriptDir string
	// where logs go
	LogDir string

	Timeout     int      `json:"timeout"`
	Interpreter string   `json:"interpreter"`
	Script      string   `json:"script"`
	Env         []string `json:"env,omitempty"`

	Command []string

	// results
	ExitCode      int   `json:"exit_code"`
	TimedOut      bool  `json:"timed_out"`
	ExecuteTimeMs int64 `json:"execute_time_ms"`
}

func (me *Exec) Init() error {

	if me.Interpreter == "" {
		me.Interpreter = "/bin/sh"
	}

	scriptfile := filepath.Join(me.ScriptDir, "job.sh")
	me.Command = []string{
		me.Interpreter,
		"-c",
		scriptfile,
	}
	slog.Debug("switching command", "cmd", me.Command)
	slog.Debug("writing script", "script", scriptfile)

	prepend := "#!/usr/bin/env bash\n"
	prepend += "set -x\n"
	if me.ImmediateExit {
		prepend += "set -e\n"
	}
	prepend += "\n"

	os.WriteFile(scriptfile, []byte(prepend+me.Script), 0755)

	return nil
}

func (me *Exec) SetEnv(k, v string) {
	kv := fmt.Sprintf("%s=%s", k, v)
	me.Env = append(me.Env, kv)
}

func (me *Exec) buildCommand(ctx context.Context) (*exec.Cmd, func(), context.Context) {
	var cancel func()

	if me.Timeout > 0 {
		slog.Debug("build command with timeout", "timeout", me.Timeout)
		ctx2, cancel2 := context.WithTimeout(ctx, time.Duration(me.Timeout)*time.Second)
		cancel = func() {
			if cancel2 == nil {
				slog.Debug("no cancel func for exec context")
			} else {
				slog.Debug("cancelling exec context")
				cancel2()
			}
		}
		ctx = ctx2
	}

	cmdline := me.Command

	slog.Debug("execute command", "cmdline", cmdline, "timeout", me.Timeout)

	c := exec.CommandContext(ctx, cmdline[0], cmdline[1:]...)
	c.Dir = me.Dir
	c.Stderr = os.Stderr
	env := make([]string, 0)
	env = append(env, os.Environ()...)
	env = append(env, me.Env...)
	c.Env = env
	slog.Debug("setting environment variable", "nenv", len(env))

	// always keep cancel func callable
	if cancel == nil {
		cancel = func() {}
	}
	return c, cancel, ctx
}

func (me *Exec) Run() error {
	slog.Debug("exec run")
	now := time.Now()

	ctx := context.Background()
	c, cancel, ctx := me.buildCommand(ctx)
	defer cancel()

	c.Stdout = me.FileLog
	c.Stderr = me.FileLog

	err := c.Run()

	dur := time.Since(now)
	me.ExecuteTimeMs = int64(dur.Nanoseconds() / 1000000)

	if ctx.Err() == context.DeadlineExceeded {
		me.TimedOut = true
	} else if err != nil {
		if c.ProcessState != nil {
			me.ExitCode = c.ProcessState.ExitCode()
		}
		return err
	}

	me.ExitCode = c.ProcessState.ExitCode()

	slog.Debug("return err", "err", err)
	return nil
}
