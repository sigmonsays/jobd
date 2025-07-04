package ui

import (
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/sigmonsays/jobd/job"
)

type IndexPage struct {
	Title string
	Jobs  []*job.JobSpec
}

func (me *Ui) Index(w http.ResponseWriter, r *http.Request) {

	// duplicate jobs array for sorting
	jobs := make([]*job.JobSpec, 0)
	for _, job := range me.Context.AppConfig.Jobs {
		jobs = append(jobs, job)
	}
	slices.SortFunc(jobs, func(a, b *job.JobSpec) int {
		return strings.Compare(a.JobId, b.JobId)
	})
	data := &IndexPage{
		Title: "jobd",
		Jobs:  jobs,
	}

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
