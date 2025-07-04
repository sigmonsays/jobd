package app

import (
	"context"
	"log/slog"

	"github.com/sigmonsays/jobd/api"
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

	return ret, nil
}
