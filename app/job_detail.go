package app

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/sigmonsays/jobd/api"
)

func (me *Api) JobDetail(context context.Context, params api.JobDetailParams) (api.JobDetailRes, error) {

	ret := &api.JobDetailResponse{}

	jid, err := strconv.ParseInt(params.Jid.Value, 10, 32)
	if err != nil {
		slog.Debug("GetConfig error", "jid", jid, "err", err)
		return &api.JobDetailInternalServerError{
			Code:    "parse-int",
			Message: err.Error(),
		}, nil
	}

	jcfg, err := me.JobConfig.GetConfig(params.Jid.Value)
	if err != nil {
		slog.Debug("GetConfig error", "jid", jid, "err", err)
		return nil, err
	}

	jres, err := me.Executor.GetResult(params.Jid.Value, 0)
	if err != nil {
		slog.Debug("GetResult error", "jid", jid, "err", err)
	}

	_ = jres
	_ = jcfg

	return ret, nil
}
