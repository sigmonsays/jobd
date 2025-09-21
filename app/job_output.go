package app

import (
	"context"
	"log/slog"

	"github.com/sigmonsays/jobd/api"
)

func (me *Api) JobOutput(context context.Context, params api.JobOutputParams) (api.JobOutputRes, error) {

	ret := &api.JobOutputResponse{}

	jid := params.Jid.Value
	jcfg, err := me.JobConfig.GetConfig(jid)
	if err != nil {
		slog.Debug("GetConfig error", "jid", jid, "err", err)
		return nil, err
	}

	jres, err := me.Executor.GetResult(jid, 0)
	if err != nil {
		slog.Debug("GetResult error", "jid", jid, "err", err)
		return nil, err
	}

	_ = jcfg

	lines, err := jres.RunSpec.GetLastLogs(1000)
	if err != nil {
		slog.Debug("GotLastLogs error", "jid", jid, "err", err)
		return nil, err
	}
	ret.Output.SetTo(lines)

	return ret, nil
}
