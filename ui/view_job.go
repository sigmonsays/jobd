package ui

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sigmonsays/jobd/job"
	"github.com/sigmonsays/jobd/schedule"
)

type ViewJobPage struct {
	Title      string
	Job        *job.JobSpec
	RunSpec    *job.RunSpec
	Incomplete bool
	NextRun    *schedule.Next
}

func (me *Ui) ViewJob(w http.ResponseWriter, r *http.Request) {
	jid := r.PathValue("jid")
	if jid == "" {
		me.handleErrorf(w, r, "missing jid")
		return
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
	jresult, err := me.Context.Executor.GetResult(jid)
	if err != nil {
	}

	incomplete := false
	if jobSpec == nil || shedJob == nil || jresult == nil {
		incomplete = true
	}

	data := &ViewJobPage{
		Title:      fmt.Sprintf("job %s", jid),
		Job:        jobSpec,
		Incomplete: incomplete,
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
