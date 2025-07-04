package app

import (
	"context"
	"log/slog"

	"github.com/sigmonsays/jobd/api"
	"github.com/sigmonsays/jobd/schedule"
)

func (me *Api) RunJob(context context.Context, req *api.RunJobRequest) (api.RunJobRes, error) {
	ret := &api.RunJobResponse{}

	// get job
	jcfg, err := me.JobConfig.GetConfig(req.Jobid.Value)
	if err != nil {
		slog.Debug("GetConfig error", "jid", req.Jobid.Value, "err", err)
		return &api.RunJobInternalServerError{
			Code:    "get-config:not-found",
			Message: err.Error(),
		}, nil
	}

	if jcfg.Disabled {
		return &api.RunJobInternalServerError{
			Code:    "job-disabled",
			Message: "job is disabled",
		}, nil
	}

	_, err = me.Scheduler.FindJobByName(req.Jobid.Value)
	if err != nil {
		slog.Debug("GetConfig error", "jid", req.Jobid.Value, "err", err)
		return &api.RunJobInternalServerError{
			Code:    "scheduler:job-not-found",
			Message: err.Error(),
		}, nil
	}

	go schedule.RunJobSpec(me.Executor, jcfg)

	return ret, nil
}
