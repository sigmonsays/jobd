package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sigmonsays/jobd/api"
	"github.com/sigmonsays/jobd/job"
	"github.com/sigmonsays/jobd/schedule"
)

func (me *Api) ListJob(context context.Context) (api.ListJobRes, error) {
	ret := &api.ListJobResponse{}

	jids := me.JobConfig.ListJobs()
	runid := 0

	for _, jid := range jids {

		jcfg, err := me.JobConfig.GetConfig(jid)
		if err != nil {
			slog.Debug("GetConfig error", "jid", jid, "err", err)
			continue
		}

		jres, err := me.Executor.GetResult(jid, runid)
		if err != nil {
			slog.Debug("GetResult error", "jid", jid, "err", err)
		}

		job := JobFromApi(jcfg, jres)

		ret.Jobs = append(ret.Jobs, *job)

	}

	return ret, nil
}

func JobFromApi(jcfg *job.JobSpec, jres *schedule.Result) *api.Job {
	job := &api.Job{}
	job.Jid.SetTo(jcfg.JobId)
	job.StepCount.SetTo(int32(len(jcfg.Steps)))

	if jres == nil {
		job.Exitcode.SetTo(-1)
	} else {
		job.Exitcode.SetTo(int32(jres.ExitCode))
		job.LatestRunid.SetTo(fmt.Sprintf("%d", jres.RunSpec.RunId))
	}
	// todo: Lots of missing fields
	// todo: jcfg.Steps
	return job
}

func RunFromApi(runSpec *job.RunSpec) *api.Run {
	ret := &api.Run{}
	ret.Runid.SetTo(fmt.Sprintf("%d", runSpec.RunId))

	m := make(map[string]string, 0)
	for k, v := range runSpec.Vars.Global.GetMap() {
		m[k] = v.String()
	}
	ret.Globals.SetTo(m)

	return ret
}
