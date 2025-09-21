package app

import (
	"context"
	"log/slog"

	"github.com/sigmonsays/jobd/api"
)

func (me *Api) JobDetail(context context.Context, params api.JobDetailParams) (api.JobDetailRes, error) {

	ret := &api.JobDetailResponse{}

	jid := params.Jid.Value
	jcfg, err := me.JobConfig.GetConfig(jid)
	if err != nil {
		slog.Debug("GetConfig error", "jid", jid, "err", err)
		return nil, err
	}

	jres, err := me.Executor.GetResult(jid, 0)
	if err != nil {
		slog.Debug("GetResult error", "jid", jid, "err", err)
	}
	job := JobFromApi(jcfg, jres)

	ret.Job.SetTo(*job)

	out := jres.RunSpec.GetLogFile()

	ret.Output.SetTo(out)

	_ = jres
	_ = jcfg

	return ret, nil
}
