package ui

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/sigmonsays/jobd/templates"
)

func (me *Ui) getTemplates() (*template.Template, error) {
	// setup functions
	fmap := map[string]any{
		"now":     me.Now,
		"version": app_version,
		"tojson":  tojson,
		"escape":  html.EscapeString,
		"nl2br":   NewLineToBr,
	}
	if me.DevMode {
		tmpl, err := template.New("jobd").
			Funcs(fmap).
			ParseGlob("./templates/*.html")
		if err != nil {
			return tmpl, err
		}
	}

	tmpl, err := template.New("jobd").
		Funcs(fmap).
		ParseFS(templates.HTML, "*.html")
	if err != nil {
		return tmpl, err
	}
	return tmpl, err
}
func NewLineToBr(x string) string {
	return strings.ReplaceAll(x, "\n", "<br>\n")
}
func app_version() string {
	// todo
	return "dev"
}

func (me *Ui) Now() string {
	return time.Now().Format(time.RFC3339)
}

func tojson(x any) (string, error) {
	v, err := json.MarshalIndent(x, "", " ")
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func (me *Ui) handleErrorf(w http.ResponseWriter, r *http.Request, s string, args ...any) {
	e := fmt.Errorf(s, args...)
	me.handleError(w, r, e)
}

func (me *Ui) handleError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Debug("http error", "error", err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	e := &AppError{
		Error: err.Error(),
	}
	buf, _ := json.Marshal(e)
	fmt.Fprintf(w, "%s\n", buf)
}

type AppError struct {
	Error string `json:"error"`
}
