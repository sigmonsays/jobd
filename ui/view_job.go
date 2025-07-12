package ui

import (
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/sigmonsays/jobd/job"
	"github.com/sigmonsays/jobd/schedule"
)

type ViewJobPage struct {
	Title      string
	Job        *job.JobSpec
	RunSpec    *job.RunSpec
	Schedule   *schedule.Job
	Incomplete bool
	NextRun    *schedule.Next
	Runs       []string
}

func (me *Ui) ViewJob(w http.ResponseWriter, r *http.Request) {
	jid := r.PathValue("jid")
	if jid == "" {
		me.handleErrorf(w, r, "missing jid")
		return
	}
	runid_str := r.PathValue("run")
	runid := 0
	if runid_str != "" {
		var err error
		runid, err = strconv.Atoi(runid_str)
		if err != nil {
			me.handleErrorf(w, r, "invalid run: %s", err)
			return
		}
	}

	// find configured job
	jobSpec, err := me.Context.JobConfig.GetConfig(jid)
	if err != nil {
		me.handleError(w, r, err)
		return
	}

	// ensure we have a scheduled job
	shedJob, err := me.Context.Scheduler.FindJobByName(jid)
	if err != nil {
	}

	// ensure we have a executing job
	jresult, err := me.Context.Executor.GetResult(jid, runid)
	if err != nil {

	}

	// list runs for job
	runbase := filepath.Join(jobSpec.Directory, "run")
	runs, err := ListRuns(runbase)
	if err != nil {

	}

	incomplete := false
	if jobSpec == nil || shedJob == nil || jresult == nil {
		incomplete = true
	}

	data := &ViewJobPage{
		Title:      fmt.Sprintf("job %s", jid),
		Job:        jobSpec,
		Schedule:   shedJob,
		Incomplete: incomplete,
		Runs:       runs,
	}
	if jresult != nil {
		data.RunSpec = jresult.RunSpec
	}

	if shedJob != nil {
		next, err := me.Context.Scheduler.FindNext(jid)
		if err == nil {
			data.NextRun = next
		}
	}

	tmpl, err := me.getTemplates()
	if err != nil {
		me.handleError(w, r, err)
		return
	}

	template_name := "job"
	if incomplete {
		template_name = "job-incomplete"
	}

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	err = tmpl.ExecuteTemplate(w, template_name, data)
	if err != nil {
		slog.Warn("Execute", "e", err)
	}
}

func ListRuns(rundir string) ([]string, error) {
	ret := make([]string, 0)
	slog.Debug("ListRuns", "rundir", rundir)

	files, err := ioutil.ReadDir(rundir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		_, err := strconv.ParseInt(file.Name(), 10, 32)
		if err != nil {
			continue
		}
		ret = append(ret, file.Name())
	}

	return ret, nil
}
