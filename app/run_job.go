package app

import (
	"context"
	"log/slog"

	"github.com/sigmonsays/jobd/api"
	"github.com/sigmonsays/jobd/job"
	"github.com/sigmonsays/jobd/schedule"
)

func (me *Api) RunJob(context context.Context, req *api.RunJobRequest) (api.RunJobRes, error) {
	ret := &api.RunJobResponse{}

	// get job
	jcfg, err := me.JobConfig.GetConfig(req.Jobid.Value)
	if err != nil {
		slog.Debug("GetConfig error", "jid", req.Jobid.Value, "err", err)
		ret.Message.SetTo(err.Error())
		return ret, nil
	}

	if jcfg.Disabled {
		ret.Message.SetTo("Job disabled")
		return ret, nil
	}

	_, err = me.Scheduler.FindJobByName(req.Jobid.Value)
	if err != nil {
		slog.Debug("GetConfig error", "jid", req.Jobid.Value, "err", err)
		ret.Message.SetTo(err.Error())
		return ret, nil
	}

	opts := schedule.DefaultExecuteOptions()
	opts.Stackable = true // always run

	runnable := &job.RunJob{
		JobSpec: jcfg,
	}

	go func() {
		err := me.Executor.Execute(req.Jobid.Value, runnable, opts)
		if err != nil {
			slog.Debug("Execute error", "jid", req.Jobid.Value, "err", err)
		}

	}()

	return ret, nil
}
