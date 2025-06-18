package ui

import (
	"log/slog"
	"net/http"

	"github.com/sigmonsays/jobd/job"
)

type IndexPage struct {
	Title string
	Jobs  []*job.JobSpec
}

func (me *Ui) Index(w http.ResponseWriter, r *http.Request) {
	data := &IndexPage{
		Title: "jobd",
		Jobs:  me.Context.AppConfig.Jobs,
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
