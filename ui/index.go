package ui

import (
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/sigmonsays/jobd/job"
	"github.com/sigmonsays/jobd/schedule"
)

type IndexPage struct {
	Title string
	Jobs  []*JobEntry
}
type JobEntry struct {
	JobSpec  *job.JobSpec
	RunSpec  *schedule.Result
	Schedule *schedule.Job

	// latest info; status string
	RunId                     int
	Status                    string
	Timestamp, HumanTimestamp string
}

func (me *Ui) Index(w http.ResponseWriter, r *http.Request) {

	// duplicate jobs array for sorting by jobid
	jobs := make([]*job.JobSpec, 0)
	for _, job := range me.Context.AppConfig.Jobs {
		jobs = append(jobs, job)
	}
	slices.SortFunc(jobs, func(a, b *job.JobSpec) int {
		return strings.Compare(a.JobId, b.JobId)
	})

	entries := make([]*JobEntry, 0)
	data := &IndexPage{
		Title: "jobd",
		Jobs:  entries,
	}

	// load job details
	for _, job := range jobs {
		jid := job.JobId

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
		jresult, err := me.Context.Executor.GetResult(jid, 0)
		jentry := &JobEntry{
			JobSpec:  jobSpec,
			RunSpec:  jresult,
			Schedule: shedJob,
		}
		if err == nil && jresult != nil {

			// go through each step to determine if the job was a success
			// todo: Fix this by saving job result in the top level RunSpec
			job_success := true
			for _, step := range jresult.RunSpec.StepResults {
				if step.ExitCode != 0 || step.Error != "" {
					job_success = false
					break
				}
			}
			if job_success {
				jentry.Status = "SUCCESS"
			} else {
				jentry.Status = "FAILED"
			}

			jentry.Timestamp = jentry.RunSpec.StartTimestamp.Format(time.RFC3339)
			jentry.HumanTimestamp = humanize.Time(jentry.RunSpec.StartTimestamp)
			jentry.RunId = jentry.RunSpec.RunSpec.RunId
		}

		entries = append(entries, jentry)
	}
	data.Jobs = entries

	tmpl, err := me.getTemplates()
	if err != nil {
		me.handleError(w, r, err)
		return
	}

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)

	err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		slog.Warn("Execute", "e", err)
	}
}
