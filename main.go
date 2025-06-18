package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	gologging "github.com/sigmonsays/go-logging"
	"github.com/sigmonsays/jobd/api"
	"github.com/sigmonsays/jobd/app"
	"github.com/sigmonsays/jobd/config"
	"github.com/sigmonsays/jobd/job"
	"github.com/sigmonsays/jobd/log"
	"github.com/sigmonsays/jobd/schedule"
	"github.com/sigmonsays/jobd/ui"
)

type Options struct {
	ConfigFile string
	LogLevel   string
}

func main() {

	opts := &Options{
		ConfigFile: "/etc/jobd.yaml",
		LogLevel:   "info",
	}
	flag.StringVar(&opts.ConfigFile, "config", opts.ConfigFile, "specify config file")
	flag.StringVar(&opts.LogLevel, "loglevel", opts.LogLevel, "log level")
	flag.Parse()
	cfg := config.GetDefaultConfig()

	fmt.Printf("Load config from %s\n", opts.ConfigFile)

	err := run(cfg, opts)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

func run(cfg *config.AppConfig, opts *Options) error {
	// setup logging
	gologging.SetLogLevel(opts.LogLevel)

	// setup slog logging
	_, mainFilename, _, _ := runtime.Caller(0)
	mainDir := filepath.Dir(mainFilename)
	slogOpts := &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source, _ := a.Value.Any().(*slog.Source)
				if source != nil {
					x, _ := filepath.Rel(mainDir, source.File)
					source.File = x
				}
			}
			return a
		},
	}
	err := log.SetLevel(slogOpts, opts.LogLevel)
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, slogOpts)).
		With("app", "jobd")
	slog.SetDefault(logger)

	loadConfig := false
	if st, err := os.Stat(opts.ConfigFile); err == nil {
		if st.IsDir() == false {
			loadConfig = true
		}
	}
	if loadConfig {
		err := cfg.LoadYaml(opts.ConfigFile)
		if err != nil {
			ExitError("LoadYaml %s: %s", opts.ConfigFile, err)
		}
	}

	// build api
	ctx := &app.Context{
		AppConfig: cfg,
	}
	app := &app.Api{
		Context: ctx,
	}
	webui := &ui.Ui{
		Context: ctx,
	}
	srv, err := api.NewServer(app)
	if err != nil {
		return err
	}

	mx := http.NewServeMux()

	mx.Handle("/api", srv)
	mx.Handle("/api/", srv)

	static := http.FileServer(http.Dir("static"))
	mx.Handle("/static/", http.StripPrefix("/static/", static))

	mx.HandleFunc("/ui", webui.Index)
	mx.HandleFunc("/ui/job/{jid}", webui.ViewJob)

	httpsrv := http.Server{
		Addr:    cfg.HttpAddr,
		Handler: mx,
	}

	// initialize jobs
	SetJobDefaults(cfg)

	jobConfig := schedule.NewJobConfig()

	// setup job directories
	for _, job := range cfg.Jobs {
		dirs := []string{
			filepath.Join(job.Directory, "run"),
		}
		for _, dir := range dirs {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				slog.Warn("Mkdir error", "dir", dir, "err", err)
			}
		}
		jobConfig.SetConfig(job.JobId, job)
	}

	// print config before starting
	cfg.PrintConfig()

	// start executor
	exec := schedule.NewExecutor()

	// start scheduler
	sched := schedule.NewScheduler()
	sched.Start()
	defer func() {
		sched.Stop()
	}()

	app.Context.Scheduler = sched
	app.Context.Executor = exec
	app.Context.JobConfig = jobConfig

	// start jobs
	for _, job_spec := range cfg.Jobs {

		rj := &job.RunJob{
			JobSpec: job_spec,
		}

		j := schedule.NewJob()
		j.Id = job_spec.JobId
		j.Schedule = job_spec.Schedule
		j.Fn = func() error {
			return exec.Execute(job_spec.JobId, rj)
		}
		sched.AddJob(j)

		if job_spec.Immediate {
			go func() {
				exec.Execute(j.Id, rj)
			}()
		}
	}

	slog.Info("starting http server", "addr", cfg.HttpAddr)
	err = httpsrv.ListenAndServe()
	if err != nil {
		ExitError("ListenAndServe %s: %s", cfg.HttpAddr, err)
	}

	return nil

}

// set defaults on the jobs
func SetJobDefaults(cfg *config.AppConfig) error {
	for _, job := range cfg.Jobs {
		if job.Directory == "" {
			job.Directory = filepath.Join(cfg.DataDir, "jobs", job.JobId)
		}
		for idx, step := range job.Steps {
			step_num := idx + 1
			if step.Shell != nil {
				step.Shell.Timeout = cfg.ShellDefaults.Timeout
			}
			if step.Id == "" {
				step.Id = fmt.Sprintf("step_%d", step_num)
			}
		}
	}
	return nil
}
